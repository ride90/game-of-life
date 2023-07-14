import { useState } from 'react';
import './App.css';
import Universe from './Universe';

function App() {
    const [universes, setUniverses] = useState([]);
    const [editing, setEditing] = useState(false);

    const createNewUniverse = () => {
        setEditing(true);
        setUniverses([...universes, { cells: Array(20).fill().map(() => Array(20).fill(false)) }]);
    };

    const saveUniverse = () => {
        setEditing(false);
    };

    return (
        <div className="App">
            {universes.map((universe, index) => (
                <div key={index} className="universe">
                    <Universe universe={universe} />
                </div>
            ))}
            <div className="buttons">
                {editing ? (
                    <button onClick={saveUniverse}>Save Universe</button>
                ) : (
                    <button onClick={createNewUniverse}>New Universe</button>
                )}
            </div>
        </div>
    );
}

export default App;
