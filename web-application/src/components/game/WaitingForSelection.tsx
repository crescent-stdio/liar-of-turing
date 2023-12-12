export default function WaitingForSelection() {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center p-4">
      <div className="bg-white shadow-xl rounded-lg w-full max-w-md mx-auto p-6">
        <h3 className="font-bold text-xl py-6 text-center">
          Players are choosing AI user...
        </h3>
      </div>
    </div>
  );
}
