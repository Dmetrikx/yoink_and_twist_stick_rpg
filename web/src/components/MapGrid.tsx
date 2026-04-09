import React from 'react';
import { GameState } from '../types/game';

interface Props {
  state: GameState;
  onMove: (dx: number, dy: number) => void;
}

export default function MapGrid({ state, onMove }: Props) {
  const { mode, player, discovered, locations, hazards, camps, cavernTiles } = state;

  let cols: number, rows: number, tileSize: number;
  if (mode === 'jungle') {
    cols = 15; rows = 15; tileSize = 40;
  } else if (mode === 'inner') {
    cols = 5; rows = 5; tileSize = 80;
  } else {
    cols = 3; rows = 8; tileSize = 60;
  }

  const innerLocations: Record<string, string> = {
    '1,1': 'Village Elder',
    '3,1': 'Shopkeeper',
    '1,3': 'Frontiersman',
    '3,3': 'Return',
    '2,1': 'Basket Weaver',
    '2,3': 'Guest Hut',
  };

  const tiles: React.ReactElement[] = [];
  for (let y = 0; y < rows; y++) {
    for (let x = 0; x < cols; x++) {
      const key = `${x},${y}`;
      const isPlayer = player.x === x && player.y === y;

      let className = 'border border-gray-700 flex items-center justify-center text-xs relative';

      if (mode === 'jungle') {
        if (!discovered[key] && !isPlayer) {
          className += ' bg-gray-800'; // fog
        } else if (locations[key]) {
          className += ' bg-green-800';
        } else if (hazards[key] && discovered[key]) {
          className += ' bg-red-900';
        } else if (camps[key]) {
          className += ' bg-yellow-800';
        } else {
          className += ' bg-green-950';
        }
      } else if (mode === 'inner') {
        if (innerLocations[key]) {
          className += ' bg-yellow-700';
        } else {
          className += ' bg-yellow-900';
        }
      } else {
        // cavern
        if (key === '1,0') {
          className += ' bg-blue-800'; // exit
        } else if (cavernTiles[key]) {
          className += ' bg-purple-900';
        } else {
          className += ' bg-gray-900';
        }
      }

      let label = '';
      if (mode === 'jungle' && discovered[key]) {
        if (locations[key]) label = locations[key].name.charAt(0);
        else if (hazards[key]) label = '!';
        else if (camps[key]) label = 'C';
      } else if (mode === 'inner') {
        if (innerLocations[key]) label = innerLocations[key].split(' ').map(w => w[0]).join('');
      } else if (mode === 'cavern') {
        if (key === '1,0') label = 'EXIT';
      }

      tiles.push(
        <div
          key={key}
          className={className}
          style={{ width: tileSize, height: tileSize }}
          title={
            mode === 'inner' ? (innerLocations[key] || '') :
            mode === 'jungle' && discovered[key] ? (locations[key]?.name || hazards[key] || (camps[key] ? 'Camp' : '')) :
            mode === 'cavern' && key === '1,0' ? 'Exit' : ''
          }
        >
          {isPlayer ? (
            <img
              src="/assets/Jungle_RPG_character.png"
              alt="Player"
              className="w-full h-full object-contain"
            />
          ) : (
            <span className="text-gray-300 select-none">{label}</span>
          )}
        </div>
      );
    }
  }

  return (
    <div>
      <div
        className="grid border border-gray-600"
        style={{
          gridTemplateColumns: `repeat(${cols}, ${tileSize}px)`,
          gridTemplateRows: `repeat(${rows}, ${tileSize}px)`,
        }}
      >
        {tiles}
      </div>
      <div className="flex gap-2 mt-3 justify-center">
        <button onClick={() => onMove(0, -1)} className="px-3 py-1 bg-gray-700 rounded hover:bg-gray-600">Up</button>
        <button onClick={() => onMove(-1, 0)} className="px-3 py-1 bg-gray-700 rounded hover:bg-gray-600">Left</button>
        <button onClick={() => onMove(0, 1)} className="px-3 py-1 bg-gray-700 rounded hover:bg-gray-600">Down</button>
        <button onClick={() => onMove(1, 0)} className="px-3 py-1 bg-gray-700 rounded hover:bg-gray-600">Right</button>
      </div>
    </div>
  );
}
