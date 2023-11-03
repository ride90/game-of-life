// Constants
const UNIVERSE_SIZE = 50;
const DEAD_CELL_COLOUR = "#2c2c2c";
const EDITABLE_CELL_COLOUR = "#434343";
const API_REQUEST_TIMEOUT = 5000;

// Next 2 lines are pretty damn sad, but I don't care tbh.
const API_URL_BASE = window.location.protocol + "//" + window.location.host + "/api";
const WS_UPDATES_URL = ((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/ws/updates";

// APIClient for making HTTP and WebSocket requests
class APIClient {
    constructor() {
        // Setup http client.
        this.initHTTP();
        // Setup WS client.
        this.initWS();
    }

    // Initialize HTTP client
    initHTTP() {
        this.axios = axios.create({
            baseURL: API_URL_BASE,
            timeout: API_REQUEST_TIMEOUT,
        });
    }

    // Initialize WebSocket client
    initWS() {
        const self = this;
        this.ws = new WebSocket(WS_UPDATES_URL);
        this.ws.onopen = () => {
            console.log("WS successfully connected.");
        };
        this.ws.onclose = (event) => {
            console.log("WS closed connection:", event);
        };
        this.ws.onerror = (error) => {
            console.log("WS error: ", error);
        };
        this.ws.onclose = function (e) {
            console.log("Socket is closed. Reconnect will be attempted in 1 second.", e.reason);
            setTimeout(function () {
                self.initWS();
            }, 1000);
        };
    }

    // Check server health
    health() {
        this.axios.get(API_URL_BASE + "/health")
            .then(function (response) {
                console.log("Health", response.statusText);
            })
            .catch(function (error) {
                console.log(error);
            });
    }

    // Create a new universe
    createUniverse(colour, cells) {
        this.axios.post(API_URL_BASE + "/universe", {
            colour: colour,
            cells: cells,
        })
            .then(function (response) {
                console.log(response);
            })
            .catch(function (error) {
                alert(error.response.data);
            });
    }

    // Reset the multiverse
    resetMultiverse() {
        this.axios.post(API_URL_BASE + "/bigbang", {})
            .then(function (response) {
                console.log(response);
            })
            .catch(function (error) {
                alert(error.response.data);
            });
    }

    // Merge universes
    mergeMultiverse() {
        this.axios.post(API_URL_BASE + "/merge", {})
            .then(function (response) {
                console.log(response);
            })
            .catch(function (error) {
                alert(error.response.data);
            });
    }
}

// Multiverse for managing universes
class Multiverse {
    constructor(universes) {
        this.universes = universes || [];
        this.isEditable = true;
        // Get API client and health ping server.
        this.apiClient = new APIClient();
        this.apiClient.health();
    }

    // Create a new universe
    createNewUniverse(isEditable) {
        let universe = new Universe(true, getRandomBrightColor());
        universe.isEditable = isEditable;
        this.universes.push(universe);
    }

    // Save the most recent universe
    saveNewUniverse() {
        // TODO: Save only not empty one.
        let universe = this.universes[this.universes.length - 1];
        // Save universe locally.
        universe.isEditable = false;
        // Create universe on the server.
        this.apiClient.createUniverse(universe.colour, universe.cells);
    }

    // Reset the multiverse
    reset() {
        this.apiClient.resetMultiverse();
    }

    // Reset the multiverse
    merge() {
        this.apiClient.mergeMultiverse();
    }

    // Remove the most recently created universe
    dropNewUniverse() {
        this.universes.pop();
    }

    // Consume updates from the WebSocket server
    consumeUpdates(onUpdate) {
        this.apiClient.ws.onmessage = (event) => {
            // Get updates from the server.
            const editableUniverses = this.universes.filter((universe) => universe.isEditable);
            this.universes = [];
            for (const data of JSON.parse(event.data)) {
                this.universes.push(new Universe(false, data.colour, data.cells));
            }
            this.universes = this.universes.concat(editableUniverses);
            // Very "Efficient" re-rendering of all non-editable universes.
            if (onUpdate) {
                onUpdate();
            }
        };
    }

    // Render the multiverse
    render() {
        let wrapper = $('<div>');
        wrapper.append(this.renderExisting());
        wrapper.append(this.renderEditable());
        return wrapper;
    }

    // Render a set of universes
    _render(universes, containerClass) {
        let table = $('<table>');
        table.addClass(containerClass);
        let row;
        for (let i = 0; i < universes.length; i++) {
            if (i % 4 === 0 || universes[i].cells[0].length > UNIVERSE_SIZE) {
                row = $('<tr>');
            }
            let td = $('<td>');
            td.append(universes[i].render());
            row.append(td);
            table.append(row);
        }
        return table;
    }

    // Render existing universes
    renderExisting() {
        return this._render(
            this.universes.filter((universe) => !universe.isEditable),
            "multiverse"
        );
    }

    // Render editable universes
    renderEditable() {
        return this._render(
            this.universes.filter((universe) => universe.isEditable),
            "wizard"
        );
    }
}

// Create an instance of Multiverse
let mu = new Multiverse([]);

// Universe class representing an individual universe
class Universe {
    constructor(isEditable, colour, cells) {
        this.isEditable = isEditable;
        this.colour = colour;
        this.cells = cells || new Array(UNIVERSE_SIZE).fill(false).map(
            () => new Array(UNIVERSE_SIZE).fill(false)
        );
    }

    // Render an editable universe
    _renderEditable() {
        let universe = this;
        let table = $('<table class="universe">');

        for (let x = 0; x < this.cells.length; x++) {
            let row = $('<tr>');
            for (let y = 0; y < this.cells[x].length; y++) {
                let td = $('<td class="cell">');
                td.attr("x", x);
                td.attr("y", y);
                // Highlight cell if it's alive.
                if (this.cells[x][y]) {
                    td.css("background", this.colour);
                }
                td.addClass("editable");
                // Cell onclick handler.
                td.on("click", () => {
                    this.cellOnClickHandler(td, universe);
                });
                // Hover highlight.
                td.hover(
                    () => this.cellOnHoverInHandler(td, universe),
                    () => this.cellOnHoverOutHandler(td, universe)
                );
                row.append(td);
            }
            table.append(row);
        }
        return table;
    }

    // Render an existing universe
    _renderExisting() {
        // Set the size of each cell and the padding between cells
        const cellSize = 6;
        const padding = 1;
        let universe = this;
        const height = (cellSize + padding) * universe.cells.length;
        const width = (cellSize + padding) * universe.cells[0].length;
        let $canvas = $('<canvas width="' + width + '" height="' + height + '">');
        let canvas = $canvas[0];
        let ctx = canvas.getContext("2d");

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
        return canvas;
    }

    // Render universe
    render() {
        if (this.isEditable) {
            return this._renderEditable();
        }
        return this._renderExisting();
    }

    // Handle cell click event
    cellOnClickHandler(td, universe) {
        let x = td.attr("x");
        let y = td.attr("y");
        universe.cells[x][y] = !universe.cells[x][y];
        if (universe.cells[x][y]) {
            td.css("background", universe.colour);
        }
    }

    // Handle cell hover in event
    cellOnHoverInHandler(td, universe) {
        if (universe.cells[td.attr("x")][td.attr("y")]) {
            return;
        }
        td.css("background", universe.colour);
    }

    // Handle cell hover out event
    cellOnHoverOutHandler(td, universe) {
        if (universe.cells[td.attr("x")][td.attr("y")]) {
            return;
        }
        td.css("background", EDITABLE_CELL_COLOUR);
    }
}

// List of predefined bright colors
const brightColors = [
    "#fa725a",
    "#FDCB58",
    "#1738ea",
    "#0099ff",
    "#ffa0ab",
    "#cd09ec",
    "#ff9507",
    "#a1ff6c",
    "#0cf632",
    "#5b09fa",
    "#cd60f8",
    "#0277fd"
];

// Global variable to keep track of the last color index
let lastColorIndex = null;

// Function to get a random index that is different from the last index
function getRandomIndex(excludeIndex, arrayLength) {
    let newIndex = Math.floor(Math.random() * arrayLength);
    while (newIndex === excludeIndex) {
        newIndex = Math.floor(Math.random() * arrayLength);
    }
    return newIndex;
}

// Function to get a random bright color that was not picked previously
function getRandomBrightColor() {
    const colorCount = brightColors.length;

    // Ensure there's more than one color to pick from
    if (colorCount <= 1) {
        return brightColors[0];
    }

    let newIndex = getRandomIndex(lastColorIndex, colorCount);

    lastColorIndex = newIndex;
    return brightColors[newIndex];
}

// Create a new universe and update UI
function newUniverse() {
    mu.createNewUniverse(true);
    $('#wizardWrapper').html(mu.renderEditable());
}

// Initialize the application
function initApp() {
    const $multiverseWrapper = $('#multiverseWrapper');
    const $wizardWrapper = $('#wizardWrapper');
    $multiverseWrapper.append(mu.renderExisting());

    // Display buttons.
    let newButton = $("#new");
    let saveButton = $("#save");
    let dropButton = $("#drop");
    let resetButton = $("#reset");
    let mergeButton = $("#merge");

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
            mu.reset();
        }
    });

    // Reset universe btn handler.
    mergeButton.on("click", () => {
        if (confirm("Make a big mess?") == true) {
            mu.merge();
        }
    });

    // Get updates and rerender them.
    mu.consumeUpdates(() => $multiverseWrapper.html(mu.renderExisting()));
}

// Initialize the app when the document is ready
$(document).ready(function () {
    initApp();
});
