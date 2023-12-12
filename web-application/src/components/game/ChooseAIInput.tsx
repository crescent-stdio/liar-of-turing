import { isFinishedSubmitionAtom } from "@/store/gameAtom";
import { User } from "@/types/playerTypes";
import { UserSelection, WsJsonRequest } from "@/types/wsTypes";
import { useAtom } from "jotai";
import { useState } from "react";

type ChooseAIInputProps = {
  userList: User[];
  userData: User;
  sendMessage: (message: WsJsonRequest) => void;
};
export default function ChooseAIInput({
  userList,
  userData,
  sendMessage,
}: ChooseAIInputProps) {
  const [AIUsername, setAIUsername] = useState<string>(userList[0].username);
  const [reason, setReason] = useState<string>("");
  const [, setIsSubmitted] = useAtom(isFinishedSubmitionAtom);

  const handleChooseAI = (event: React.FormEvent<HTMLFormElement>) => {
    console.log(AIUsername, reason);
    event.preventDefault();
    const user_selection: UserSelection = {
      user: userData,
      selection: AIUsername,
      reason: reason,
    };
    const jsonData: WsJsonRequest = {
      max_player: 0,
      action: "choose_ai",
      user: userData,
      timestamp: Date.now(),
      message: reason,
      game_round: 0,
      game_turns_left: 0,
      user_selection: user_selection,
    };
    sendMessage(jsonData);
    setIsSubmitted(true);
  };
  // const handleChange;
  return (
    <>
      <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center p-4">
        <div className="bg-white shadow-xl rounded-lg w-full max-w-md mx-auto p-6">
          <h3 className="font-bold text-lg mb-4 text-center">Choose AI user</h3>
          <form className="flex flex-col space-y-4" onSubmit={handleChooseAI}>
            <div className="flex flex-col">
              <label htmlFor="ai" className="text-sm text-gray-600 mb-2">
                I think the AI is...
              </label>
              <select
                name="ai"
                id="ai"
                className="border border-gray-300 rounded-md p-2 focus:border-liar-blue focus:ring-liar-blue"
                onChange={(e) => setAIUsername(e.target.value)}
              >
                {userList &&
                  userList.length > 0 &&
                  userList.map((user: User, index) => {
                    return (
                      <option key={index} value={user.username}>
                        {user.username}
                      </option>
                    );
                  })}
              </select>
            </div>

            <div className="flex flex-col">
              <label htmlFor="reason" className="text-sm text-gray-600 mb-2">
                Reason
              </label>
              <input
                className="border border-gray-300 rounded-md p-2 focus:border-liar-blue focus:ring-liar-blue"
                type="text"
                id="reason"
                value={reason}
                onChange={(e) => setReason(e.target.value)}
              />
            </div>

            <button
              type="submit"
              className="bg-liar-blue hover:bg-liar-blue-dark text-white font-bold py-2 px-4 rounded"
            >
              Submit
            </button>
          </form>
        </div>
      </div>
    </>
  );
}