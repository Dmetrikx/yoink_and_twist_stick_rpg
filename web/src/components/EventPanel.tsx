import { GameState } from '../types/game';

interface Props {
  state: GameState;
}

const locationImages: Record<string, string> = {
  'Capybara Den': '/assets/Capybara_Den.png',
  'Ruins': '/assets/Ruins.png',
  'Waterfall': '/assets/Waterfall.png',
  'Cavern': '/assets/Cavern.png',
  'Village': '/assets/Purchase_Village.png',
};

const hazardImages: Record<string, string> = {
  'Beehive': '/assets/Beehive.png',
  'Jaguar': '/assets/Jaguar.png',
  'Anaconda': '/assets/Anaconda.png',
  'Quicksand': '/assets/Quicksand.png',
};

export default function EventPanel({ state }: Props) {
  const { player, mode, locations, hazards, lastEvent } = state;
  const key = `${player.x},${player.y}`;

  let imageSrc: string | null = null;

  if (mode === 'jungle') {
    const loc = locations[key];
    if (loc && locationImages[loc.name]) {
      imageSrc = locationImages[loc.name];
    }
    const haz = hazards[key];
    if (haz && hazardImages[haz]) {
      imageSrc = hazardImages[haz];
    }
  }

  return (
    <div className="bg-gray-800 rounded-lg p-3">
      <h3 className="font-bold mb-2 text-amber-400">Events</h3>
      {imageSrc && (
        <img
          src={imageSrc}
          alt="Location"
          className="w-full max-w-[300px] rounded mb-2 mx-auto"
        />
      )}
      <p className="text-sm">{lastEvent || 'Your adventure awaits...'}</p>
    </div>
  );
}
