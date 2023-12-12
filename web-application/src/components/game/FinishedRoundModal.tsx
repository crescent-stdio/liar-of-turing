import { useAtomValue } from "jotai";
import ChooseAIInput from "./ChooseAIInput";
import WaitingForSelection from "./WaitingForSelection";
import { isFinishedRoundAtom, isFinishedSubmitionAtom } from "@/store/gameAtom";
import { playerListAtom, userAtom } from "@/store/chatAtom";
import { WsJsonRequest } from "@/types/wsTypes";
type FinishedRoundModalProps = {
  sendMessage: (message: WsJsonRequest) => void;
};
export default function FinishedRoundModal({
  sendMessage: handleWebSocketMessageSend,
}: FinishedRoundModalProps) {
  const user = useAtomValue(userAtom);
  const playerList = useAtomValue(playerListAtom);
  const isFinishedRound = useAtomValue(isFinishedRoundAtom);
  const isFinishedSubmition = useAtomValue(isFinishedSubmitionAtom);

  if (!user || !playerList) return <></>;
  if (!isFinishedRound) return <></>;
  return (
    <>
      {user.player_type === "player" ? (
        !isFinishedSubmition ? (
          <ChooseAIInput
            userData={user}
            userList={playerList}
            sendMessage={handleWebSocketMessageSend}
          />
        ) : (
          <WaitingForSelection />
        )
      ) : (
        <WaitingForSelection />
      )}
    </>
  );
}
