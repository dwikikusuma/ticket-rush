import grpc
from concurrent import futures
import time
import random
import sys
import os
import redis
import pricing_pb2
import pricing_pb2_grpc
from grpc_reflection.v1alpha import reflection

sys.path.append(os.path.dirname(os.path.abspath(__file__)))

class PricingService(pricing_pb2_grpc.PricingServiceServicer):
    
    def __init__(self):
        self.r = redis.Redis(host='localhost', port=6379, db=0, decode_responses=True)
        print("🔌 Connected to Redis")

    def GetPrice(self, request, context):
        event_id = request.event_id
        
        view_count_key = f"view_count:{event_id}"
        current_views = self.r.incr(view_count_key)
        
        print(f"💰 Calculating price for Event: {event_id} (Views: {current_views})")

        if current_views > 10:
            print(f"   🔥 High demand detected! ({current_views} views)")
            multiplier = 1.5
        
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