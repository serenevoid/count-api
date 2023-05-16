import requests
import concurrent.futures
import time
import csv

URL = "http://192.168.1.5:8080"

def send_req(url):
    start_time = time.time()
    _ = requests.get(url)
    end_time = time.time()
    return end_time - start_time

def test_api(num_requests):
    response_times = []

    with concurrent.futures.ThreadPoolExecutor() as executor:
        # Send multiple requests concurrently
        futures = [executor.submit(send_req, URL) for _ in range(num_requests)]

        # Collect the response times
        for future in concurrent.futures.as_completed(futures):
            response_time = future.result()
            response_times.append(response_time)

    return response_times

num_req = 10000000

start_time = time.time()
response_times = test_api(num_req)
end_time = time.time()

print(f'Total requests sent: {num_req}')
print(f'Total time taken: {end_time - start_time} seconds')
print(f'Average response time: {sum(response_times) / len(response_times)} seconds')
print(f'Minimum response time: {min(response_times)} seconds')
print(f'Maximum response time: {max(response_times)} seconds')

# Save the response times to a CSV file
with open('response_times.csv', 'w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(['Request Index', 'Response Time (seconds)'])
    writer.writerows(enumerate(response_times))

print("Response times saved to response_times.csv")
