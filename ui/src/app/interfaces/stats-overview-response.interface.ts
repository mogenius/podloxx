export interface IStatsOverviewResponse {
  packetsSum: number;
  transmitBytes: number;
  receivedBytes: number;
  unknownBytes: number;
  localTransmitBytes: number;
  localReceivedBytes: number;
  totalNodes: number;
  totalNamespaces: number;
  totalPods: number;
  uptime: string;
}
