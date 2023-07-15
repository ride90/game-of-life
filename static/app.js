const UNIVERSE_SIZE = 20;
const DEAD_CELL_COLOUR = "#2c2c2c";


class Multiverse {
    constructor(universes) {
        this.universes = universes || []
    }

    createEmptyUniverse(isEditable) {
        let universe = new Universe(true, getRandomColor());
        universe.isEditable = isEditable;
        this.universes.push(universe)
    }

    render() {
        let table = $('<table class="multiverse">');
        let row
        for (let i = 0; i < this.universes.length; i++) {

            if (i % 5 === 0) {
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


class Cell {
    constructor(alive, universe) {
        this.alive = alive;
        this.universe = universe;
    }
}

class Universe {
    constructor(isEditable, colour, cells) {
        this.isEditable = isEditable;
        this.colour = colour;
        this.cells = cells || Array(UNIVERSE_SIZE).fill(Array(UNIVERSE_SIZE).fill(new Cell(false, this)));
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

                // Cell click.
                td.on("click", () => {
                    console.log(this.cells);
                    this.cellOnClickHandler(td, universe)
                    console.log(this.cells);
                    // this.cells[i][j] = !this.cells[i][j]
                });
                // Hover highlight.
                // td.hover(
                //     () => this.cellOnHoverInHandler(td, universe),
                //     () => this.cellOnHoverOutHandler(td, universe),
                // )
                // // Highlight live cell.
                // if (this.cells[i][j]) {
                //     td.toggleClass("live")
                //     td.css("background", this.colour)
                // }
                row.append(td);
            }
            table.append(row);
        }
        return table;
    }

    cellOnClickHandler(td, universe) {
        td.toggleClass("live")
        if (td.hasClass("live")) {
            td.css("background", universe.colour)
        }
    }

    cellOnHoverInHandler(cell, universe) {
        if (cell.hasClass("live")) {
            return
        }
        cell.css("background", universe.colour)
    }

    cellOnHoverOutHandler(cell, universe) {
        if (cell.hasClass("live")) {
            return
        }
        cell.css("background", DEAD_CELL_COLOUR)
    }
}

function getRandomColor() {
    let letters = '0123456789ABCDEF';
    let color = '#';

    while (true) {
        for (let i = 0; i < 6; i++) {
            color += letters[Math.floor(Math.random() * 16)];
        }

        // Check if the generated color is not too dark or too light
        let luminance = getLuminance(color);
        if (luminance > 0.1 && luminance < 0.9) {
            break;
        }
        color = '#'; // Reset the color if it doesn't meet the criteria
    }

    return color;
}

function getLuminance(color) {
    // Extract RGB values from the color string
    let r = parseInt(color.substr(1, 2), 16);
    let g = parseInt(color.substr(3, 2), 16);
    let b = parseInt(color.substr(5, 2), 16);

    // Calculate relative luminance using the formula for sRGB color space
    return (0.2126 * r + 0.7152 * g + 0.0722 * b) / 255;
}

function newUniverse() {
    mu.createEmptyUniverse(true)
    $('#multiverseWrapper').html(mu.render());
}

function initApp() {
    $('#multiverseWrapper').append(mu.render());

    // Display buttons.

    let newButton = $("#new")
    newButton.show();
    newButton.on("click", newUniverse);

}

$(document).ready(function () {
    initApp();
});

