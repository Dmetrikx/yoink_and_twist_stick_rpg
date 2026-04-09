import { Player } from '../types/game';

interface Props {
  player: Player;
  basketCap: number;
}

export default function StatsBar({ player, basketCap }: Props) {
  return (
    <div className="flex gap-4 p-3 bg-gray-800 rounded-lg text-sm">
      <StatBar label="HP" value={player.hp} max={100} color="bg-red-500" />
      <StatBar label="Energy" value={player.energy} max={100} color="bg-blue-500" />
      <div className="flex items-center gap-2">
        <span className="text-yellow-400 font-bold">Gold:</span>
        <span className="text-yellow-300">{player.gold} / {basketCap}</span>
      </div>
      <div className="flex items-center gap-2">
        <span className="text-gray-400">Pos:</span>
        <span>({player.x}, {player.y})</span>
      </div>
    </div>
  );
}

function StatBar({ label, value, max, color }: { label: string; value: number; max: number; color: string }) {
  const pct = Math.max(0, Math.min(100, (value / max) * 100));
  return (
    <div className="flex items-center gap-2 min-w-[140px]">
      <span className="font-bold w-14">{label}:</span>
      <div className="flex-1 bg-gray-700 rounded-full h-4 relative overflow-hidden">
        <div className={`${color} h-full rounded-full transition-all`} style={{ width: `${pct}%` }} />
        <span className="absolute inset-0 flex items-center justify-center text-xs text-white font-bold">
          {value}
        </span>
      </div>
    </div>
  );
}
