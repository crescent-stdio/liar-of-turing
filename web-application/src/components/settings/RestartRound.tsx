import { adminUser, initialUserSelection } from "@/store/chatStore";
import { WsJsonRequest } from "@/types/wsTypes";

type RestartRoundProps = {
  sendMessage: (message: WsJsonRequest) => void;
};
export default function RestartRound({ sendMessage }: RestartRoundProps) {
  const handleRestartRound = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    const jsonData: WsJsonRequest = {
      max_player: 0,
      action: "restart_round",
      user: adminUser,
      message: `Restart Round`,
      timestamp: Date.now(),
      game_round: 0,
      game_turns_left: 0,
      user_selection: initialUserSelection,
    };
    sendMessage(jsonData);
  };

  return (
    <div>
      <h3 className="mt-4 mb-2 font-bold text-xl">Restart Round</h3>

      <button
        onClick={handleRestartRound}
        className="px-2 py-1 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-liar-blue focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75"
      >
        Restart Round
      </button>
    </div>
  );
}
