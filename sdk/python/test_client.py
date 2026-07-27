from client import SwarmLensClient
client =SwarmLensClient(brokers="localhost:19092")
client.emit(
    swarm_id="swarm-python-test-1",
    agent_id="python-agent",
    content="hello from the python sdk",
    recipient_agent_id="go-agent",
)
print("event published successfully")