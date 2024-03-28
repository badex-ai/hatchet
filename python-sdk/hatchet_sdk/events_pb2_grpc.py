# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from . import events_pb2 as events__pb2

class EventsServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Push = channel.unary_unary(
                '/EventsService/Push',
                request_serializer=events__pb2.PushEventRequest.SerializeToString,
                response_deserializer=events__pb2.Event.FromString,
                )
        self.ReplaySingleEvent = channel.unary_unary(
                '/EventsService/ReplaySingleEvent',
                request_serializer=events__pb2.ReplayEventRequest.SerializeToString,
                response_deserializer=events__pb2.Event.FromString,
                )
        self.PutLog = channel.unary_unary(
                '/EventsService/PutLog',
                request_serializer=events__pb2.PutLogRequest.SerializeToString,
                response_deserializer=events__pb2.PutLogResponse.FromString,
                )


class EventsServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Push(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ReplaySingleEvent(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def PutLog(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_EventsServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Push': grpc.unary_unary_rpc_method_handler(
                    servicer.Push,
                    request_deserializer=events__pb2.PushEventRequest.FromString,
                    response_serializer=events__pb2.Event.SerializeToString,
            ),
            'ReplaySingleEvent': grpc.unary_unary_rpc_method_handler(
                    servicer.ReplaySingleEvent,
                    request_deserializer=events__pb2.ReplayEventRequest.FromString,
                    response_serializer=events__pb2.Event.SerializeToString,
            ),
            'PutLog': grpc.unary_unary_rpc_method_handler(
                    servicer.PutLog,
                    request_deserializer=events__pb2.PutLogRequest.FromString,
                    response_serializer=events__pb2.PutLogResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'EventsService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class EventsService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Push(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/EventsService/Push',
            events__pb2.PushEventRequest.SerializeToString,
            events__pb2.Event.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ReplaySingleEvent(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/EventsService/ReplaySingleEvent',
            events__pb2.ReplayEventRequest.SerializeToString,
            events__pb2.Event.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def PutLog(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/EventsService/PutLog',
            events__pb2.PutLogRequest.SerializeToString,
            events__pb2.PutLogResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
