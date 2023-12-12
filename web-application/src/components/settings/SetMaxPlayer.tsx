import { adminUser } from "@/store/chatStore";
import { maxPlayerAtom } from "@/store/gameAtom";
import { WsJsonRequest } from "@/types/wsTypes";
import { useAtom } from "jotai";
import { useState } from "react";
type SetMaxPayerProps = {
  sendMessage: (message: WsJsonRequest) => void;
};
export default function SetMaxPayer({ sendMessage }: SetMaxPayerProps) {
  const [maxPlayer, setMaxPlayer] = useAtom(maxPlayerAtom);
  const [player, setPlayer] = useState<number>(1);

  const handleSetMaxPlayer = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (player === maxPlayer) return;
    const jsonData: WsJsonRequest = {
      max_player: player,
      action: "set_max_player",
      timestamp: Date.now(),
      user: adminUser,
      message: `Set max player to ${player}`,
      game_round: 0,
      game_turns_left: 0,
    };
    sendMessage(jsonData);
    setMaxPlayer(player);
  };

  return (
    <div>
      <h3 className="mt-4 mb-2 font-bold text-xl">
        Set max player - current max player is{" "}
        <span className="text-[#3b82f6]">{maxPlayer}</span>
      </h3>
      <form className="flex flex-row" onSubmit={handleSetMaxPlayer}>
        <label htmlFor="username">Username</label>
        <select
          name="username"
          id="username"
          onChange={(e) => {
            setPlayer(+e.target.value);
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
          Set Max Player
        </button>
      </form>
    </div>
  );
}
