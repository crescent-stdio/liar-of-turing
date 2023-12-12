type SetGameStatusProps = {
  isGameStarted: boolean;
  isYourTurn: boolean;
};
export default function SetGameStatus() {
  return (
    <div>
      <div className="w-full border-t border-gray-300 my-4"></div>
      <div className="hidden lg:block h-full border-l border-gray-300"></div>
    </div>
  );
}
