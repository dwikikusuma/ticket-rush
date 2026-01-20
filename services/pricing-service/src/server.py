import grpc
from concurrent import futures
import time
import random
import sys
import os

import pricing_pb2
import pricing_pb2_grpc
from grpc_reflection.v1alpha import reflection

sys.path.append(os.path.dirname(os.path.abspath(__file__)))

class PricingService(pricing_pb2_grpc.PricingServiceServicer):
    
    def GetPrice(self, request: pricing_pb2.PriceRequest, context) -> pricing_pb2.PriceResponse:
        
        print(f"💰 Calculating price for: {request.event_id} (Seat: {request.seat_id})")
        
        is_surge = random.random() < 0.3
        multiplier = 1.5 if is_surge else 1.0
        
        original_price = request.base_price
        final_price = int(original_price * multiplier)

        return pricing_pb2.PriceResponse(
            final_price=final_price,
            multiplier=multiplier
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pricing_pb2_grpc.add_PricingServiceServicer_to_server(PricingService(), server)

    SERVICE_NAMES = (
        "pricing.PricingService",
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    port = "50051"
    server.add_insecure_port('[::]:' + port)
    print(f"🐍 Python Pricing Service running on port {port}")
    
    server.start()
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        print("stopping pricing service")
        server.stop(0)

if __name__ == '__main__':
    serve()