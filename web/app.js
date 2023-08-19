const UNIVERSE_SIZE = 50;
const DEAD_CELL_COLOUR = "#2c2c2c";
const EDITABLE_CELL_COLOUR = "#434343";
const API_URL_BASE = "http://127.0.0.1:4000/api"
const API_REQUEST_TIMEOUT = 5000
const WS_UPDATES_URL = "ws://127.0.0.1:4000/ws/updates"


class APIClient {
    constructor() {
        // Setup http client.
        this.initHTTP()
        // Setup WS client.
        this.initWS()
    }

    initHTTP() {
        this.axios = axios.create({
            baseURL: API_URL_BASE,
            timeout: API_REQUEST_TIMEOUT
        });
    }

    initWS() {
        const self = this;
        this.ws = new WebSocket(WS_UPDATES_URL);
        this.ws.onopen = () => {
            console.log("WS successfully connected.");
        };
        this.ws.onclose = event => {
            console.log("WS closed connection:", event);
        };
        this.ws.onerror = error => {
            console.log("WS error: ", error);
        };
        this.ws.onclose = function (e) {
            console.log("Socket is closed. Reconnect will be attempted in 1 second.", e.reason);
            setTimeout(function () {
                self.initWS();
            }, 1000);
        };
    }

    health() {
        this.axios.get(API_URL_BASE + "/health")
            .then(function (response) {
                console.log("Health", response.statusText);
            })
            .catch(function (error) {
                console.log(error);
            })
    }

    createUniverse(colour, cells) {
        this.axios.post(API_URL_BASE + "/universe", {
            colour: colour,
            cells: cells
        })
            .then(function (response) {
                console.log(response);
            })
            .catch(function (error) {
                alert(error.response.data);
            });
    }

    resetMultiverse() {
        this.axios.post(API_URL_BASE + "/bigbang", {})
            .then(function (response) {
                console.log(response);
            })
            .catch(function (error) {
                alert(error.response.data);
            });
    }
}


class Multiverse {
    constructor(universes) {
        this.universes = universes || []
        this.isEditable = true
        // Get API client and health ping server.
        this.apiClient = new APIClient()
        this.apiClient.health()
    }

    createNewUniverse(isEditable) {
        let universe = new Universe(true, getRandomColor());
        universe.isEditable = isEditable;
        this.universes.push(universe)
    }

    saveNewUniverse() {
        // TODO: Save only not empty one.
        let universe = this.universes[this.universes.length - 1]
        // Save universe locally.
        universe.isEditable = false
        // Create universe on the server.
        this.apiClient.createUniverse(universe.colour, universe.cells)

    }

    reset() {
        this.apiClient.resetMultiverse()
    }

    dropNewUniverse() {
        this.universes.pop();
    }

    consumeUpdates(onUpdate) {
        this.apiClient.ws.onmessage = event => {
            // Get updates from the server.
            const editableUniverses = this.universes.filter((universe) => universe.isEditable);
            this.universes = []
            for (const data of JSON.parse(event.data)) {
                this.universes.push(new Universe(false, data.colour, data.cells))
            }
            this.universes = this.universes.concat(editableUniverses)
            // Very "Efficient" re-rendering of all non-editable universes.
            if (onUpdate) {
                onUpdate()
            }
        }
    }

    render() {
        let wrapper = $('<div>');
        wrapper.append(this.renderExisting())
        wrapper.append(this.renderEditable())
        return wrapper
    }

    _render(universes, containerClass) {
        let table = $('<table>');
        table.addClass(containerClass)
        let row
        for (let i = 0; i < universes.length; i++) {
            if (i % 4 === 0) {
                row = $('<tr>');
            }
            let td = $('<td>');
            td.append(universes[i].render())
            row.append(td);
            table.append(row);
        }
        return table;
    }

    renderExisting() {
        return this._render(
            this.universes.filter((universe) => !universe.isEditable),
            "multiverse"
        );
    }

    renderEditable() {
        return this._render(
            this.universes.filter((universe) => universe.isEditable),
            "wizard"
        );
    }
}

let mu = new Multiverse([])


class Universe {
    constructor(isEditable, colour, cells) {
        this.isEditable = isEditable;
        this.colour = colour;
        this.cells = cells || new Array(UNIVERSE_SIZE).fill(false).map(
            () => new Array(UNIVERSE_SIZE).fill(false)
        );
    }

    _renderEditable() {
        let universe = this
        let table = $('<table class="universe">');

        for (let x = 0; x < this.cells.length; x++) {
            let row = $('<tr>');
            for (let y = 0; y < this.cells[x].length; y++) {
                let td = $('<td class="cell">');
                td.attr("x", x);
                td.attr("y", y);
                // Highlight cell if it's alive.
                if (this.cells[x][y]) {
                    td.css("background", this.colour)
                }
                td.addClass("editable")
                // Cell onclick handler.
                td.on("click", () => {
                    this.cellOnClickHandler(td, universe)
                });
                // Hover highlight.
                td.hover(
                    () => this.cellOnHoverInHandler(td, universe),
                    () => this.cellOnHoverOutHandler(td, universe),
                )
                row.append(td);
            }
            table.append(row);
        }
        return table;
    }

    _renderExisting() {
        // Set the size of each cell and the padding between cells
        const cellSize = 6;
        const padding = 1;
        const size = (cellSize + padding) * UNIVERSE_SIZE
        let universe = this
        let $canvas = $('<canvas width="' + size + '" height="' + size + '">');
        let canvas = $canvas[0]
        let ctx = canvas.getContext("2d")

        // Loop through the matrix and render each cell
        for (let row = 0; row < universe.cells.length; row++) {
            for (let col = 0; col < universe.cells[row].length; col++) {
                const cellValue = universe.cells[row][col];
                // Calculate the position for the current cell
                const x = col * (cellSize + padding);
                const y = row * (cellSize + padding);
                // Set the fill color based on the cell value
                ctx.fillStyle = cellValue === true ? universe.colour : DEAD_CELL_COLOUR;
                // Draw the cell
                ctx.fillRect(x, y, cellSize, cellSize);
            }
        }
        return canvas
    }

    render() {
        if (this.isEditable) {
            return this._renderEditable()
        }
        return this._renderExisting()
    }

    cellOnClickHandler(td, universe) {
        let x = td.attr("x")
        let y = td.attr("y")
        universe.cells[x][y] = !universe.cells[x][y];
        if (universe.cells[x][y]) {
            td.css("background", universe.colour)
        }
    }

    cellOnHoverInHandler(td, universe) {
        if (universe.cells[td.attr("x")][td.attr("y")]) {
            return
        }
        td.css("background", universe.colour)
    }

    cellOnHoverOutHandler(td, universe) {
        if (universe.cells[td.attr("x")][td.attr("y")]) {
            return
        }
        td.css("background", EDITABLE_CELL_COLOUR)
    }
}

let lastColor = null;

function colorDistance(color1, color2) {
    return Math.sqrt(
        Math.pow(parseInt(color1.substr(1, 2), 16) - parseInt(color2.substr(1, 2), 16), 2) +
        Math.pow(parseInt(color1.substr(3, 2), 16) - parseInt(color2.substr(3, 2), 16), 2) +
        Math.pow(parseInt(color1.substr(5, 2), 16) - parseInt(color2.substr(5, 2), 16), 2)
    );
}

function getRandomColor() {
    let letters = '0123456789ABCDEF';
    let color = '#';

    while (true) {
        for (let i = 0; i < 6; i++) {
            color += letters[Math.floor(Math.random() * 16)];
        }

        let luminance = getLuminance(color);
        if (luminance > 0.6 && luminance < 0.9) {
            // Check that the color is different enough from the last color
            if (lastColor === null || colorDistance(color, lastColor) > 100) {
                break;
            }
        }
        color = '#'; // Reset the color if it doesn't meet the criteria
    }

    lastColor = color;
    return color;
}

function getLuminance(color) {
    let r = parseInt(color.substr(1, 2), 16);
    let g = parseInt(color.substr(3, 2), 16);
    let b = parseInt(color.substr(5, 2), 16);

    return (0.2126 * r + 0.7152 * g + 0.0722 * b) / 255;
}

function newUniverse() {
    mu.createNewUniverse(true)
    $('#wizardWrapper').html(mu.renderEditable());
}

function initApp() {
    const $multiverseWrapper = $('#multiverseWrapper')
    const $wizardWrapper = $('#wizardWrapper')
    $multiverseWrapper.append(mu.renderExisting());

    // Display buttons.
    let newButton = $("#new")
    let saveButton = $("#save")
    let dropButton = $("#drop")
    let resetButton = $("#reset")
    // New universe btn handler.
    newButton.show();
    newButton.on("click", () => {
        newUniverse();
        newButton.hide();
        saveButton.show();
        dropButton.show();
    });
    // Save universe btn handler.
    saveButton.on("click", () => {
        mu.saveNewUniverse();
        $wizardWrapper.html(mu.renderEditable());
        newButton.show();
        saveButton.hide();
        dropButton.hide();
    });
    // Drop universe btn handler.
    dropButton.on("click", () => {
        mu.dropNewUniverse();
        $wizardWrapper.html(mu.renderEditable());
        newButton.show();
        saveButton.hide();
        dropButton.hide();
    });
    // Reset universe btn handler.
    resetButton.on("click", () => {
        if (confirm("Are you sure you want to destroy everything?") == true) {
            mu.reset()
        }
    });

    // Get updates and rerender them.
    mu.consumeUpdates(
        () => $multiverseWrapper.html(mu.renderExisting())
    );
}

$(document).ready(function () {
    initApp();
});

