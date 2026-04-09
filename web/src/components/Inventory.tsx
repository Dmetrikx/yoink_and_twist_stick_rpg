import { InventorySlot } from '../types/game';

interface Props {
  inventory: (InventorySlot | null)[];
  onUse: (slot: number) => void;
}

export default function Inventory({ inventory, onUse }: Props) {
  return (
    <div className="bg-gray-800 rounded-lg p-3">
      <h3 className="font-bold mb-2 text-amber-400">Inventory</h3>
      <div className="flex gap-2">
        {[0, 1, 2].map(i => {
          const item = inventory[i];
          return (
            <div key={i} className="border border-gray-600 rounded p-2 w-28 text-center bg-gray-900">
              {item ? (
                <>
                  <div className="text-sm font-semibold">{item.name}</div>
                  <div className="text-xs text-gray-400">x{item.count}</div>
                  <button
                    onClick={() => onUse(i)}
                    className="mt-1 px-2 py-0.5 bg-green-700 rounded text-xs hover:bg-green-600"
                  >
                    Use
                  </button>
                </>
              ) : (
                <div className="text-gray-600 text-sm">Empty</div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
