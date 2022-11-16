export interface IStatsOverviewResponse {
  externalBandwidthPerSec: string;
  internalBandwidthPerSec: string;
  lastUpdate: number;
  localReceivedBytes: number;
  localTransmitBytes: number;
  packetsPerSec: number;
  packetsSum: number;
  receivedBytes: number;
  totalNamespaces: number;
  totalNodes: number;
  totalPods: number;
  transmitBytes: number;
  unknownBytes: number;
  uptime: string;
}
