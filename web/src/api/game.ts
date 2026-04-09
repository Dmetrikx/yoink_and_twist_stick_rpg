import { GameState } from '../types/game';
import { apiFetch } from './client';

export function newGame(): Promise<GameState> {
  return apiFetch<GameState>('/api/game/new', { method: 'POST' });
}

export function getState(): Promise<GameState> {
  return apiFetch<GameState>('/api/game/state');
}

export function move(dx: number, dy: number): Promise<GameState> {
  return apiFetch<GameState>('/api/game/move', {
    method: 'POST',
    body: JSON.stringify({ dx, dy }),
  });
}

export function doAction(action: string): Promise<GameState> {
  return apiFetch<GameState>('/api/game/action', {
    method: 'POST',
    body: JSON.stringify({ action }),
  });
}

export function useItem(slot: number): Promise<GameState> {
  return apiFetch<GameState>('/api/game/item', {
    method: 'POST',
    body: JSON.stringify({ slot }),
  });
}
