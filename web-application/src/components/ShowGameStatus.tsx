import { playerListAtom } from "@/store/chatAtom";
import {
  gameRoundAtom,
  gameRoundNumAtom,
  gameTurnsLeftAtom,
  gameTurnsNumAtom,
} from "@/store/gameAtom";
import { useAtomValue } from "jotai";

export default function ShowGameStatus() {
  const playerList = useAtomValue(playerListAtom);
  const gameRound = useAtomValue(gameRoundAtom);
  const gameTurnsLeft = useAtomValue(gameTurnsLeftAtom);
  const gameRoundNum = useAtomValue(gameRoundNumAtom);
  const gameTurnsNum = useAtomValue(gameTurnsNumAtom);
  return (
    // <div className="mx-auto w-[80vw] max-w-2xl min-h-screen relative">
    <div className="absolute top-0 -right-4 md:-right-24">
      <div className="flex flex-col shadow-md rounded-lg overflow-hidden xl:text-xl ">
        <div className="px-4 py-2 border-b border-gray-200">
          <div className="text-sm text-gray-600">Round:</div>
          <div className="font-bold text-gray-900">
            {gameRound}/{gameRoundNum}
          </div>
        </div>
        <div className="px-4 py-2">
          <div className="text-sm text-gray-600">Turns left:</div>
          <div className="font-bold text-gray-900">
            {gameTurnsLeft}/{gameTurnsNum * playerList.length}
          </div>
        </div>
      </div>
    </div>
    // </div>
  );
}
