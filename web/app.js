const UNIVERSE_SIZE = 20;
const DEAD_CELL_COLOUR = "#2c2c2c";
const EDITABLE_CELL_COLOUR = "#434343";
const API_URL_BASE = "http://127.0.0.1:4000/api"
const API_REQUEST_TIMEOUT = 5000


class APIClient {
    constructor() {
        this.axios = axios.create({
            baseURL: API_URL_BASE,
            timeout: API_REQUEST_TIMEOUT
        });

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
                console.log(error);
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
        let universe = this.universes[this.universes.length - 1]
        // Save universe locally.
        universe.isEditable = false
        // Create universe on the server.
        this.apiClient.createUniverse(universe.colour, universe.cells)

    }

    dropNewUniverse() {
        this.universes.pop();
    }

    render() {
        let table = $('<table class="multiverse">');
        let row
        for (let i = 0; i < this.universes.length; i++) {

            if (i % 4 === 0) {
                row = $('<tr>');
            }
            let td = $('<td>');
            td.append(this.universes[i].render())
            row.append(td);
            table.append(row);
        }
        return table;
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

    render() {
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
                if (this.isEditable) {
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
                }
                row.append(td);
            }
            table.append(row);
        }
        return table;
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
        if (luminance > 0.1 && luminance < 0.9) {
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
    $('#multiverseWrapper').html(mu.render());
}

function initApp() {
    $('#multiverseWrapper').append(mu.render());

    // Display buttons.
    let newButton = $("#new")
    let saveButton = $("#save")
    let dropButton = $("#drop")
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
        $('#multiverseWrapper').html(mu.render());
        newButton.show();
        saveButton.hide();
        dropButton.hide();
    });
    // Drop universe btn handler.
    dropButton.on("click", () => {
        mu.dropNewUniverse();
        $('#multiverseWrapper').html(mu.render());
        newButton.show();
        saveButton.hide();
        dropButton.hide();
    });

}

$(document).ready(function () {
    initApp();
});

