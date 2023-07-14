import { useState } from 'react';
import './Universe.css';

function Universe({ universe }) {
    const [cells, setCells] = useState(universe.cells);

    const toggleCellState = (rowIndex, cellIndex) => {
        const newCells = [...cells];
        newCells[rowIndex][cellIndex] = !newCells[rowIndex][cellIndex];
        setCells(newCells);
    };

    return (
        <table className="universe-table">
            <tbody>
            {cells.map((row, rowIndex) => (
                <tr key={rowIndex}>
                    {row.map((cell, cellIndex) => (
                        <td
                            key={cellIndex}
                            className={cell ? 'cell live' : 'cell'}
                            onClick={() => toggleCellState(rowIndex, cellIndex)}
                        />
                    ))}
                </tr>
            ))}
            </tbody>
        </table>
    );
}

export default Universe;
