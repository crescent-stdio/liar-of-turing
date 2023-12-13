export default function VerticalLine() {
  return (
    <div className="flex justify-center items-center">
      {/* Vertical Divider for larger screens */}
      <div className="hidden lg:block h-full border-l border-gray-300"></div>

      {/* Horizontal Divider for smaller screens */}
      <div className="lg:hidden w-full border-t border-gray-300 my-2"></div>
    </div>
  );
}
