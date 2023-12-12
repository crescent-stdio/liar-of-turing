import { gameRoundAtom, gameTurnsLeftAtom } from "@/store/gameAtom";
import { useAtomValue } from "jotai";

export default function ShowGameStatus() {
  const gameRound = useAtomValue(gameRoundAtom);
  const gameTurnsLeft = useAtomValue(gameTurnsLeftAtom);
  return (
    <>
      <div className="absolute top-0 right-0 mt-8 mr-8">
        <div className="flex flex-col bg-white shadow-md rounded-lg overflow-hidden">
          <div className="px-4 py-2 border-b border-gray-200">
            <div className="text-sm text-gray-600">Round:</div>
            <div className="font-bold text-gray-900">{gameRound}</div>
          </div>
          <div className="px-4 py-2">
            <div className="text-sm text-gray-600">Turns left:</div>
            <div className="font-bold text-gray-900">{gameTurnsLeft}</div>
          </div>
        </div>
      </div>
    </>
  );
}
