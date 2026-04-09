export interface Player {
  x: number; y: number; hp: number; energy: number; gold: number;
}

export interface InventorySlot {
  name: string; count: number;
}

export interface Quest {
  target: string; done: boolean;
}

export interface Location {
  name: string; enter?: boolean;
}

export interface Camp {
  uses: number;
}

export type Mode = "jungle" | "inner" | "cavern";

export interface GameState {
  player: Player;
  mode: Mode;
  discovered: Record<string, boolean>;
  locations: Record<string, Location>;
  hazards: Record<string, string>;
  camps: Record<string, Camp>;
  inventory: (InventorySlot | null)[];
  quests: Quest[];
  basketCap: number;
  cavernTiles: Record<string, string>;
  lastEvent: string;
  causeOfDeath: string;
  gameOver: boolean;
  won: boolean;
}

export interface SaveMeta {
  ID: number;
  UserID: string;
  Name: string;
  SavedAt: string;
}
