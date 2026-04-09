import { useCallback, useEffect, useState } from 'react';
import { GameState } from '../types/game';
import * as gameApi from '../api/game';
import * as savesApi from '../api/saves';
import MapGrid from '../components/MapGrid';
import StatsBar from '../components/StatsBar';
import Inventory from '../components/Inventory';
import QuestLog from '../components/QuestLog';
import EventPanel from '../components/EventPanel';
import SaveLoadModal from '../components/SaveLoadModal';

export default function GamePage() {
  const [state, setState] = useState<GameState | null>(null);
  const [showSaves, setShowSaves] = useState(false);
  const [error, setError] = useState('');

  // Try to load existing game, if none start new
  useEffect(() => {
    gameApi.getState()
      .then(setState)
      .catch(() => {
        gameApi.newGame().then(setState).catch(err => setError(err.message));
      });
  }, []);

  const handleMove = useCallback(async (dx: number, dy: number) => {
    if (!state || state.gameOver || state.won) return;
    try {
      const s = await gameApi.move(dx, dy);
      setState(s);
    } catch (err: any) {
      setError(err.message);
    }
  }, [state]);

  const handleAction = useCallback(async (action: string) => {
    if (!state || state.gameOver || state.won) return;
    try {
      const s = await gameApi.doAction(action);
      setState(s);
    } catch (err: any) {
      setError(err.message);
    }
  }, [state]);

  const handleUseItem = useCallback(async (slot: number) => {
    if (!state || state.gameOver || state.won) return;
    try {
      const s = await gameApi.useItem(slot);
      setState(s);
    } catch (err: any) {
      setError(err.message);
    }
  }, [state]);

  const handleNewGame = async () => {
    try {
      const s = await gameApi.newGame();
      setState(s);
      setError('');
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleLoadSave = async (id: number) => {
    try {
      const s = await savesApi.loadSave(id);
      setState(s);
      setShowSaves(false);
    } catch (err: any) {
      setError(err.message);
    }
  };

  // Keyboard controls
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (showSaves) return;
      switch (e.key) {
        case 'ArrowUp': e.preventDefault(); handleMove(0, -1); break;
        case 'ArrowDown': e.preventDefault(); handleMove(0, 1); break;
        case 'ArrowLeft': e.preventDefault(); handleMove(-1, 0); break;
        case 'ArrowRight': e.preventDefault(); handleMove(1, 0); break;
        case 'c': handleAction('buildCamp'); break;
        case 'e': {
          if (!state) break;
          const key = `${state.player.x},${state.player.y}`;
          if (state.mode === 'jungle') {
            const loc = state.locations[key];
            if (loc?.name === 'Village') handleAction('enterVillage');
            else if (loc?.name === 'Cavern') handleAction('enterCavern');
            else if (state.camps[key]) handleAction('rest');
          } else if (state.mode === 'inner') {
            const innerLocs: Record<string, string> = {
              '1,1': 'Village Elder', '3,1': 'Shopkeeper', '1,3': 'Frontiersman',
              '3,3': 'Return', '2,1': 'Basket Weaver', '2,3': 'Guest Hut',
            };
            const name = innerLocs[key];
            if (name === 'Guest Hut') handleAction('rest');
            else if (name === 'Basket Weaver') handleAction('basketWeave');
          }
          break;
        }
        case 'w': {
          if (!state) break;
          if (state.mode === 'inner') handleAction('exitVillage');
          else if (state.mode === 'cavern') handleAction('exitCavern');
          else if (state.mode === 'jungle') {
            const key = `${state.player.x},${state.player.y}`;
            const loc = state.locations[key];
            if (loc?.enter) handleAction('enterVillage');
          }
          break;
        }
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [handleMove, handleAction, state, showSaves]);

  if (!state) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          {error ? (
            <div>
              <p className="text-red-400 mb-4">{error}</p>
              <a href="/auth/google" className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-500">
                Sign in with Google
              </a>
            </div>
          ) : (
            <p className="text-gray-400">Loading...</p>
          )}
        </div>
      </div>
    );
  }

  // Contextual action buttons based on current tile
  const currentKey = `${state.player.x},${state.player.y}`;
  const currentLoc = state.locations[currentKey];
  const innerLocs: Record<string, string> = {
    '1,1': 'Village Elder', '3,1': 'Shopkeeper', '1,3': 'Frontiersman',
    '3,3': 'Return', '2,1': 'Basket Weaver', '2,3': 'Guest Hut',
  };
  const currentInner = state.mode === 'inner' ? innerLocs[currentKey] : null;

  return (
    <div className="min-h-screen p-4">
      <div className="max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-2xl font-bold text-amber-400">Jungle RPG</h1>
          <div className="flex gap-2">
            <button
              onClick={() => setShowSaves(true)}
              className="px-3 py-1 bg-gray-700 rounded text-sm hover:bg-gray-600"
            >
              Save/Load
            </button>
            <button
              onClick={handleNewGame}
              className="px-3 py-1 bg-gray-700 rounded text-sm hover:bg-gray-600"
            >
              New Game
            </button>
          </div>
        </div>

        <StatsBar player={state.player} basketCap={state.basketCap} />

        {(state.gameOver || state.won) && (
          <div className={`mt-4 p-4 rounded-lg text-center text-lg font-bold ${
            state.won ? 'bg-green-900 text-green-200' : 'bg-red-900 text-red-200'
          }`}>
            {state.won ? 'You purchased the village! You Win!' : `Game Over: ${state.causeOfDeath}`}
            <button
              onClick={handleNewGame}
              className="ml-4 px-4 py-1 bg-gray-700 rounded text-sm hover:bg-gray-600"
            >
              Play Again
            </button>
          </div>
        )}

        <div className="mt-4 grid grid-cols-1 lg:grid-cols-[auto_1fr] gap-4">
          <div>
            <MapGrid state={state} onMove={handleMove} />

            {/* Context actions */}
            <div className="flex flex-wrap gap-2 mt-3">
              {state.mode === 'jungle' && currentLoc?.name === 'Village' && (
                <button onClick={() => handleAction('enterVillage')} className="px-3 py-1 bg-green-700 rounded text-sm hover:bg-green-600">Enter Village (E)</button>
              )}
              {state.mode === 'jungle' && currentLoc?.name === 'Cavern' && (
                <button onClick={() => handleAction('enterCavern')} className="px-3 py-1 bg-purple-700 rounded text-sm hover:bg-purple-600">Enter Cavern (E)</button>
              )}
              {state.mode === 'jungle' && state.camps[currentKey] && state.camps[currentKey].uses > 0 && (
                <button onClick={() => handleAction('rest')} className="px-3 py-1 bg-yellow-700 rounded text-sm hover:bg-yellow-600">Rest at Camp (E)</button>
              )}
              {state.mode === 'jungle' && !state.gameOver && !state.won && (
                <button onClick={() => handleAction('buildCamp')} className="px-3 py-1 bg-yellow-800 rounded text-sm hover:bg-yellow-700">Build Camp (C) - 50g</button>
              )}
              {state.mode === 'inner' && currentInner === 'Village Elder' && (
                <button onClick={() => handleAction('addQuest')} className="px-3 py-1 bg-blue-700 rounded text-sm hover:bg-blue-600">Get Quest</button>
              )}
              {state.mode === 'inner' && currentInner === 'Shopkeeper' && (
                <>
                  <button onClick={() => handleAction('buyPotion')} className="px-3 py-1 bg-red-700 rounded text-sm hover:bg-red-600">Health Potion - 50g</button>
                  <button onClick={() => handleAction('buyStim')} className="px-3 py-1 bg-blue-700 rounded text-sm hover:bg-blue-600">Stimulant - 25g</button>
                </>
              )}
              {state.mode === 'inner' && currentInner === 'Basket Weaver' && (
                <button onClick={() => handleAction('basketWeave')} className="px-3 py-1 bg-yellow-700 rounded text-sm hover:bg-yellow-600">Weave Baskets (E)</button>
              )}
              {state.mode === 'inner' && currentInner === 'Guest Hut' && (
                <button onClick={() => handleAction('rest')} className="px-3 py-1 bg-green-700 rounded text-sm hover:bg-green-600">Rest (E)</button>
              )}
              {state.mode === 'inner' && currentInner === 'Frontiersman' && (
                <button onClick={() => handleAction('buyVillage')} className="px-3 py-1 bg-amber-600 rounded text-sm hover:bg-amber-500">Buy Village - 1500g</button>
              )}
              {state.mode === 'inner' && (currentInner === 'Return' || !currentInner) && (
                <button onClick={() => handleAction('exitVillage')} className="px-3 py-1 bg-gray-700 rounded text-sm hover:bg-gray-600">Leave Village (W)</button>
              )}
              {state.mode === 'cavern' && (
                <button onClick={() => handleAction('exitCavern')} className="px-3 py-1 bg-gray-700 rounded text-sm hover:bg-gray-600">Exit Cavern (W)</button>
              )}
            </div>
          </div>

          <div className="space-y-4">
            <EventPanel state={state} />
            <Inventory inventory={state.inventory} onUse={handleUseItem} />
            <QuestLog quests={state.quests} />
          </div>
        </div>

        {error && (
          <div className="mt-4 p-2 bg-red-900 rounded text-red-200 text-sm">
            {error}
          </div>
        )}
      </div>

      <SaveLoadModal
        open={showSaves}
        onClose={() => setShowSaves(false)}
        onLoad={handleLoadSave}
      />
    </div>
  );
}
