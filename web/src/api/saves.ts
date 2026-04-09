import { GameState, SaveMeta } from '../types/game';
import { apiFetch } from './client';

export function listSaves(): Promise<SaveMeta[]> {
  return apiFetch<SaveMeta[]>('/api/saves/');
}

export function createSave(name: string): Promise<void> {
  return apiFetch<void>('/api/saves/', {
    method: 'POST',
    body: JSON.stringify({ name }),
  });
}

export function updateSave(id: number, name: string): Promise<void> {
  return apiFetch<void>(`/api/saves/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ name }),
  });
}

export function loadSave(id: number): Promise<GameState> {
  return apiFetch<GameState>(`/api/saves/${id}/load`);
}

export function deleteSave(id: number): Promise<void> {
  return apiFetch<void>(`/api/saves/${id}`, { method: 'DELETE' });
}
