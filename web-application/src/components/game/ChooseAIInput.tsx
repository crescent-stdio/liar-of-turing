import { isFinishedSubmitionAtom } from "@/store/gameAtom";
import { SELECTION_OPEN_TIME } from "@/store/gameStore";
import { User } from "@/types/playerTypes";
import { UserSelection, WsJsonRequest } from "@/types/wsTypes";
import { useAtom } from "jotai";
import { useEffect, useState } from "react";

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
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

  const handleChooseAI = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const user_selection: UserSelection = {
      timestamp: Date.now(),
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
      game_round_num: 0,
      game_turn_num: 0,
      user_selection: user_selection,
    };
    sendMessage(jsonData);
    setIsSubmitted(true);
  };
  // Toggle modal on click
  const toggleModal = () => {
    setIsModalOpen((prev) => !prev);
  };
  const handleModalClick = (event: React.MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();
  };

  // Set up global click listener
  useEffect(() => {
    const timer = setTimeout(() => {
      setIsModalOpen(true);
    }, SELECTION_OPEN_TIME);

    window.addEventListener("click", toggleModal);

    return () => {
      clearTimeout(timer);
      window.removeEventListener("click", toggleModal);
    };
  }, []);

  return (
    <div className="cursor-pointer flex justify-center items-center h-screen">
      {isModalOpen && (
        <div className="bg-black bg-opacity-50 fixed inset-0 p-4 flex justify-center items-center h-full">
          <div
            className="bg-white shadow-xl rounded-lg w-full max-w-md mx-auto p-6"
            onClick={handleModalClick}
          >
            <h3 className="font-bold text-lg mb-4 text-center">
              Choose AI user
            </h3>
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
                    userList.map((user, index) => (
                      <option key={index} value={user.username}>
                        {user.username}
                      </option>
                    ))}
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
            <p className="text-sm text-gray-600 mt-4 text-center">
              Click other place to close the modal
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
