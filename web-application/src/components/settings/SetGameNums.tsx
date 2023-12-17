import { adminUser, initialUserSelection } from "@/store/chatStore";
import { gameRoundAtom, gameTurnsNumAtom } from "@/store/gameAtom";
import { WsJsonRequest } from "@/types/wsTypes";
import { useAtom } from "jotai";
import { useState } from "react";

type SetGameNumsProps = {
  sendMessage: (message: WsJsonRequest) => void;
};
export default function SetGameNums({ sendMessage }: SetGameNumsProps) {
  const [gameRound, setGameRound] = useAtom(gameRoundAtom);
  const [gameTurnsNum, setGameTurnsNum] = useAtom(gameTurnsNumAtom);
  const [nowGameRound, setNowGameRound] = useState(gameRound);
  const [nowGameTurnsNum, setNowGameTurnsNum] = useState(gameTurnsNum);
  const [gameRoundNum, setGameRoundNum] = useAtom(gameRoundAtom);
  const [gameTurnsLeft, setGameTurnsLeft] = useAtom(gameTurnsNumAtom);
  const handleSetGameRound = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (nowGameRound <= gameRound) return;
    const jsonData: WsJsonRequest = {
      max_player: 0,
      action: "set_game_round",
      user: adminUser,
      message: `Set Game Round`,
      timestamp: Date.now(),
      game_round: nowGameRound,
      game_turns_left: gameTurnsLeft,
      game_turn_num: gameTurnsNum,
      game_round_num: gameRoundNum,
      user_selection: initialUserSelection,
    };
    sendMessage(jsonData);
    setGameRound(nowGameRound);
  };
  const handleSetGameTurnsNum = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (nowGameTurnsNum <= gameTurnsNum) return;
    const jsonData: WsJsonRequest = {
      max_player: 0,
      action: "set_game_turn",
      user: adminUser,
      message: `Set Game Turn`,
      timestamp: Date.now(),
      game_round: gameRound,
      game_round_num: gameRoundNum,
      game_turns_left: gameTurnsLeft,
      game_turn_num: nowGameTurnsNum,
      user_selection: initialUserSelection,
    };
    sendMessage(jsonData);
    setGameTurnsNum(nowGameTurnsNum);
  };
  return (
    <div>
      <div className="font-bold">
        <h3 className="mt-4 mb-2 font-bold text-xl">Set Game Nums</h3>
        {gameRound && (
          <p className="mb-2">
            Current Game Round:{" "}
            <span className="text-[#3b82f6]">{gameRound}</span>
          </p>
        )}
        {gameTurnsNum && (
          <p className="mb-2">
            Current Game Turns Num:{" "}
            <span className="text-[#3b82f6]">{gameTurnsNum}</span>
          </p>
        )}
      </div>
      <form className="mb-2" onSubmit={handleSetGameRound}>
        <label htmlFor="gameRound">Set Game Round</label>
        <select
          name="gameRound"
          id="gameRound"
          onChange={(e) => {
            setNowGameRound(+e.target.value);
          }}
          className="border-2 border-gray-400 rounded-md w-fit-content mx-1"
        >
          {Array(10)
            .fill(0)
            .map((_, index) => {
              return (
                <option key={index} value={index + 1}>
                  {index + 1}
                </option>
              );
            })}
        </select>
        <button
          type="submit"
          className="px-2 py-1 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-liar-blue focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75"
        >
          Set Game rounds
        </button>
      </form>

      <form className="mb-2" onSubmit={handleSetGameTurnsNum}>
        <label htmlFor="gameTurnsNum">Set Game Turns Num</label>
        <select
          name="gameTurnsNum"
          id="gameTurnsNum"
          onChange={(e) => {
            setNowGameTurnsNum(+e.target.value);
          }}
          className="border-2 border-gray-400 rounded-md w-fit-content mx-1"
        >
          {Array(10)
            .fill(0)
            .map((_, index) => {
              return (
                <option key={index} value={index + 1}>
                  {index + 1}
                </option>
              );
            })}
        </select>
        <button
          type="submit"
          className="px-2 py-1 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-[#3b82f6] focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75"
        >
          Set Game Turns Num
        </button>
      </form>
    </div>
  );
}
