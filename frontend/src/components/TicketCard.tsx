import { Ticket } from "../types";

interface Props {
  ticket: Ticket;
}

export default function TicketCard({ ticket }: Props) {
  // Format price to currency (e.g. $1,200)
  const formattedPrice = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    maximumFractionDigits: 0,
  }).format(ticket.price);

  return (
    <div className="border border-gray-700 bg-gray-900 rounded-lg p-4 shadow-lg hover:border-blue-500 transition-colors">
      <div className="flex justify-between items-start">
        <div>
          <h3 className="text-xl font-bold text-white">{ticket.event_name}</h3>
          <p className="text-gray-400 text-sm">{ticket.stadium}</p>
        </div>
        <span className="bg-blue-900 text-blue-200 text-xs px-2 py-1 rounded-full">
          {ticket.status}
        </span>
      </div>

      <div className="mt-4 flex justify-between items-end">
        <div className="text-gray-500 text-sm">Seat: <span className="text-white">{ticket.seat_id}</span></div>
        <div className="text-green-400 text-2xl font-bold">{formattedPrice}</div>
      </div>
    </div>
  );
}