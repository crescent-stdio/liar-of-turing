import { SELECTION_OPEN_TIME } from "@/store/gameStore";
import { useEffect, useState } from "react";

export default function WaitingForSelection() {
  const [isModalOpen, setIsModalOpen] = useState(false);
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
    <div className="cursor-pointer flex justify-center items-center h-scree">
      {isModalOpen && (
        <div className="bg-black bg-opacity-50 fixed inset-0 p-4 flex justify-center items-center h-full">
          <div
            className="bg-white shadow-xl rounded-lg w-full max-w-md mx-auto p-6"
            onClick={handleModalClick}
          >
            <h3 className="font-bold text-xl py-6 text-center">
              Players are choosing AI user...
            </h3>
            <p className="text-sm text-gray-600 mt-4 text-center">
              Click other place to close the modal
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
