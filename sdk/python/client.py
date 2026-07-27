import json
import time
import uuid
from typing import Optional

from kafka import KafkaProducer


class SwarmLensClient:
    """Emits agent events to a SwarmLens Kafka topic, matching the Go schema."""

    def __init__(self, brokers, topic="agent.messages"):
        self.topic = topic
        self.producer = KafkaProducer(
            bootstrap_servers=brokers,
            value_serializer=lambda v: json.dumps(v).encode("utf-8"),
            key_serializer=lambda k: k.encode("utf-8"),
        )

    def emit(
        self,
        swarm_id: str,
        agent_id: str,
        content: str,
        recipient_agent_id: Optional[str] = None,
        parent_task_id: Optional[str] = None,
        framework: str = "python-sdk",
    ):
        event = {
            "event_id": str(uuid.uuid4()),
            "swarm_id": swarm_id,
            "agent_id": agent_id,
            "event_type": "message",
            "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
            "payload": {
                "content": content,
            },
            "metadata": {
                "framework": framework,
                "sdk_version": "0.1.0",
            },
        }

        if recipient_agent_id is not None:
            event["payload"]["recipient_agent_id"] = recipient_agent_id
        if parent_task_id is not None:
            event["parent_task_id"] = parent_task_id

        self.producer.send(self.topic, key=swarm_id, value=event)
        self.producer.flush()