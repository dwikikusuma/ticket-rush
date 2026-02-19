export default function SkeletonCard() {
  return (
    <div className="border border-gray-800 bg-gray-900 rounded-lg p-4 animate-pulse">
      <div className="flex justify-between items-start">
        <div className="space-y-2 flex-1">
          <div className="h-5 bg-gray-700 rounded w-3/4" />
          <div className="h-3 bg-gray-800 rounded w-1/2" />
        </div>
        <div className="h-5 bg-gray-700 rounded-full w-20 ml-4" />
      </div>
      <div className="mt-4 flex justify-between items-end">
        <div className="h-3 bg-gray-800 rounded w-24" />
        <div className="h-7 bg-gray-700 rounded w-16" />
      </div>
    </div>
  );
}