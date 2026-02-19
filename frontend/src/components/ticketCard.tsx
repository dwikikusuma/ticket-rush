import type { Ticket } from "@/types";
import StatusBadge from "@/components/statusBadge";
import { formatPrice } from "@/lib/utils";

interface Props {
  ticket: Ticket;
  onBook?: (ticket: Ticket) => void;
  isLoggedIn?: boolean;
}

export default function TicketCard({ ticket, onBook, isLoggedIn }: Props) {
  return (
    <div className="border border-gray-700 bg-gray-900 rounded-lg p-4 shadow-lg hover:border-blue-500 transition-colors">
      <div className="flex justify-between items-start gap-2">
        <div className="min-w-0">
          <h3 className="text-base font-bold text-white truncate">{ticket.event_name}</h3>
          <p className="text-gray-400 text-sm truncate">{ticket.stadium}</p>
        </div>
        <StatusBadge status={ticket.status} />
      </div>

      <div className="mt-4 flex justify-between items-end">
        <div className="text-gray-500 text-sm">
          Seat: <span className="text-white">{ticket.seat_id}</span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-green-400 text-xl font-bold">
            {formatPrice(ticket.price)}
          </span>
          {ticket.status === "AVAILABLE" && onBook && (
            <button
              onClick={() => onBook(ticket)}
              className={`px-3 py-1 rounded text-sm font-medium transition
                ${isLoggedIn
                  ? "bg-blue-600 hover:bg-blue-500 text-white"
                  : "bg-gray-700 text-gray-400 cursor-not-allowed"
                }`}
              title={isLoggedIn ? "Book ticket" : "Log in to book"}
            >
              Book
            </button>
          )}
        </div>
      </div>
    </div>
  );
}