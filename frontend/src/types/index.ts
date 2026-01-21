export interface Ticket {
  id: number;
  event_name: string;
  stadium: string;
  price: number;
  seat_id: string;
  status: string;
}

export interface SearchResponse {
  data: Ticket[];
  next_cursor: string;
}