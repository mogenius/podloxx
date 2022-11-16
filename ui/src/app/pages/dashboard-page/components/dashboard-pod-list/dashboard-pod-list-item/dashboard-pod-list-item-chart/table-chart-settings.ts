import { ChartConfiguration } from 'chart.js';

export const tableChartConfig: ChartConfiguration = {
  type: 'line',
  data: {
    datasets: [
      {
        data: [],
        fill: true,
        backgroundColor: 'rgba(255, 255, 255, 1)'
      }
    ],
    labels: []
  },
  options: {
    aspectRatio: 4,
    elements: {
      line: {
        tension: 0,
        borderWidth: 1.5
      },
      point: {
        radius: 0,
        hitRadius: 20,
        hoverRadius: 1
      }
    },
    scales: {
      'y-axis-0': {
        display: false,
        beginAtZero: true,
        suggestedMin: 0
      },
      'x-axis-0': {
        display: false
      }
    },
    responsive: true,
    plugins: {
      legend: {
        display: false
      },
      title: {
        display: false
      },
      tooltip: {
        backgroundColor: '#00495d',
        titleColor: 'white',
        bodyColor: 'white',
        displayColors: false,
        callbacks: {
          title: (context) => {
            return `${context[0].dataset.label ?? 'n/a'}`;
          },
          label: (context) => {
            return `${context.label} - ${formatBytes(+context.parsed.y, 2)}`;
          }
        }
      }
    }
  }
};

const formatBytes = (val: number | string, precision: number) => {
  let label = 'MB';
  let parseCount = 0;
  val = +val;
  while (val > 1024) {
    val = val / 1024;
    parseCount++;
  }
  switch (parseCount) {
    case 0:
      label = 'Bytes';
      break;
    case 1:
      label = 'KB';
      break;
    case 2:
      label = 'MB';
      break;
    case 3:
      label = 'GB';
      break;
    case 4:
      label = 'TB';
      break;
    case 5:
      label = 'PB';
      break;
  }

  return `${val.toFixed(precision)} ${label}`;
};
