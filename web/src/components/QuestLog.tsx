import { Quest } from '../types/game';

interface Props {
  quests: Quest[];
}

export default function QuestLog({ quests }: Props) {
  return (
    <div className="bg-gray-800 rounded-lg p-3">
      <h3 className="font-bold mb-2 text-amber-400">Quests</h3>
      {quests.length === 0 ? (
        <div className="text-gray-500 text-sm">No quests. Visit the Village Elder.</div>
      ) : (
        <ul className="space-y-1 text-sm">
          {quests.map((q, i) => (
            <li key={i} className={q.done ? 'line-through text-gray-500' : ''}>
              {q.done ? '\u2713' : '\u25CB'} Find the {q.target}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
