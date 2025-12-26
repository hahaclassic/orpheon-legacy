import React from 'react';

interface Segment {
  idx: number;
  totalStreams: number;
  range: [number, number];
}

interface TrackStatsChartProps {
  duration: number;
  progress: number;
}

// Моковые данные для графика
const mockSegments: Segment[] = [
  { idx: 0, totalStreams: 10, range: [0, 10] },
  { idx: 1, totalStreams: 30, range: [10, 20] },
  { idx: 2, totalStreams: 50, range: [20, 30] },
  { idx: 3, totalStreams: 20, range: [30, 40] },
  { idx: 4, totalStreams: 60, range: [40, 50] },
  { idx: 5, totalStreams: 40, range: [50, 60] },
  { idx: 6, totalStreams: 80, range: [60, 70] },
  { idx: 7, totalStreams: 30, range: [70, 80] },
  { idx: 8, totalStreams: 20, range: [80, 90] },
  { idx: 9, totalStreams: 10, range: [90, 100] },
];

const WIDTH = 1000; // Теперь график всегда на всю ширину контейнера
const HEIGHT = 38;
const PADDING = 8;
const LINE_COLOR = '#ff4081';
const FILL_COLOR = 'rgba(255,64,129,0.10)';

const getMaxStreams = (segments: Segment[]) =>
  Math.max(...segments.map(s => s.totalStreams), 1);

const TrackStatsChart: React.FC<TrackStatsChartProps> = () => {
  const maxStreams = getMaxStreams(mockSegments);
  const points = mockSegments.map((seg, i) => {
    const x = (i / (mockSegments.length - 1)) * (WIDTH - 2 * PADDING) + PADDING;
    const y = HEIGHT - PADDING - (seg.totalStreams / maxStreams) * (HEIGHT - 2 * PADDING);
    return [x, y];
  });

  // Area path (вся заливка)
  const areaPath = points.reduce((acc, point, i, arr) => {
    if (i === 0) return `M ${point[0]},${HEIGHT - PADDING} L ${point[0]},${point[1]}`;
    const prev = arr[i - 1];
    const cpx = (prev[0] + point[0]) / 2;
    return acc + ` C ${cpx},${prev[1]} ${cpx},${point[1]} ${point[0]},${point[1]}`;
  }, '');
  const areaPathFull = areaPath + ` L ${points[points.length - 1][0]},${HEIGHT - PADDING} Z`;

  // Линия графика
  const linePath = points.reduce((acc, point, i, arr) => {
    if (i === 0) return `M ${point[0]},${point[1]}`;
    const prev = arr[i - 1];
    const cpx = (prev[0] + point[0]) / 2;
    return acc + ` C ${cpx},${prev[1]} ${cpx},${point[1]} ${point[0]},${point[1]}`;
  }, '');

  return (
    <svg width="100%" height={HEIGHT} viewBox={`0 0 ${WIDTH} ${HEIGHT}`} style={{ display: 'block', background: 'none' }}>
      {/* Заливка под графиком */}
      <path d={areaPathFull} fill={FILL_COLOR} />
      {/* Основная линия */}
      <path d={linePath} fill="none" stroke={LINE_COLOR} strokeWidth={3} strokeLinecap="round" />
    </svg>
  );
};

export default TrackStatsChart; 