import { useEffect, useState } from 'react';
import { SaveMeta } from '../types/game';
import * as savesApi from '../api/saves';

interface Props {
  open: boolean;
  onClose: () => void;
  onLoad: (id: number) => void;
}

export default function SaveLoadModal({ open, onClose, onLoad }: Props) {
  const [saves, setSaves] = useState<SaveMeta[]>([]);
  const [newName, setNewName] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (open) {
      savesApi.listSaves().then(setSaves).catch(() => {});
    }
  }, [open]);

  if (!open) return null;

  const handleCreate = async () => {
    const name = newName.trim() || `Save ${new Date().toLocaleString()}`;
    setLoading(true);
    await savesApi.createSave(name);
    setNewName('');
    const updated = await savesApi.listSaves();
    setSaves(updated);
    setLoading(false);
  };

  const handleOverwrite = async (id: number, name: string) => {
    setLoading(true);
    await savesApi.updateSave(id, name);
    const updated = await savesApi.listSaves();
    setSaves(updated);
    setLoading(false);
  };

  const handleDelete = async (id: number) => {
    setLoading(true);
    await savesApi.deleteSave(id);
    const updated = await savesApi.listSaves();
    setSaves(updated);
    setLoading(false);
  };

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg p-5 w-[400px] max-h-[80vh] overflow-auto">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-bold text-amber-400">Save / Load</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-white text-xl">&times;</button>
        </div>

        <div className="mb-4 flex gap-2">
          <input
            type="text"
            placeholder="Save name..."
            value={newName}
            onChange={e => setNewName(e.target.value)}
            className="flex-1 bg-gray-700 rounded px-2 py-1 text-sm"
          />
          <button
            onClick={handleCreate}
            disabled={loading}
            className="px-3 py-1 bg-green-700 rounded text-sm hover:bg-green-600 disabled:opacity-50"
          >
            New Save
          </button>
        </div>

        {saves.length === 0 ? (
          <p className="text-gray-500 text-sm">No saves yet.</p>
        ) : (
          <ul className="space-y-2">
            {saves.map(s => (
              <li key={s.ID} className="bg-gray-900 rounded p-2 flex items-center justify-between">
                <div>
                  <div className="text-sm font-semibold">{s.Name}</div>
                  <div className="text-xs text-gray-500">{new Date(s.SavedAt).toLocaleString()}</div>
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => onLoad(s.ID)}
                    className="px-2 py-0.5 bg-blue-700 rounded text-xs hover:bg-blue-600"
                  >
                    Load
                  </button>
                  <button
                    onClick={() => handleOverwrite(s.ID, s.Name)}
                    className="px-2 py-0.5 bg-yellow-700 rounded text-xs hover:bg-yellow-600"
                  >
                    Overwrite
                  </button>
                  <button
                    onClick={() => handleDelete(s.ID)}
                    className="px-2 py-0.5 bg-red-700 rounded text-xs hover:bg-red-600"
                  >
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
