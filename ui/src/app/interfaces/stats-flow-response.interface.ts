export interface IStatsFlowResponse {
  [key: string]: {
    packetsSum: number;
    transmitBytes: number;
    receivedBytes: number;
    unknownBytes: number;
    localTransmitBytes: number;
    localReceivedBytes: number;
  };
}
