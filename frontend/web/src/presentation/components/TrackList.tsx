import { List, Divider } from '@mui/material';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import type { DragEndEvent } from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import TrackItem from './TrackItem';
import type { Track } from '../types';

interface TrackListProps {
  tracks: Track[];
  onAddToPlaylist?: (event: React.MouseEvent<HTMLElement>, track: Track) => void;
  showTrackNumber?: boolean;
  showAlbumLink?: boolean;
  onTrackReorder?: (sourceIndex: number, destinationIndex: number) => void;
  isDraggable?: boolean;
  onTrackClick?: (trackId: string) => void;
}

interface SortableTrackItemProps {
  track: Track;
  tracks: Track[];
  index: number;
  onAddToPlaylist?: (event: React.MouseEvent<HTMLElement>, track: Track) => void;
  showTrackNumber?: boolean;
  showAlbumLink?: boolean;
  onTrackClick?: (trackId: string) => void;
}

const SortableTrackItem = ({
  track,
  tracks,
  index,
  onAddToPlaylist,
  showTrackNumber,
  showAlbumLink,
  onTrackClick,
}: SortableTrackItemProps) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: track.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    backgroundColor: isDragging ? 'rgba(0, 0, 0, 0.1)' : 'transparent',
    cursor: 'grab',
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      <TrackItem
        track={track}
        tracks={tracks}
        index={index}
        onAddToPlaylist={onAddToPlaylist}
        showTrackNumber={showTrackNumber}
        showAlbumLink={showAlbumLink}
        onTrackClick={onTrackClick}
      />
    </div>
  );
};

const TrackList = ({
  tracks,
  onAddToPlaylist,
  showTrackNumber = true,
  showAlbumLink = true,
  onTrackReorder,
  isDraggable = false,
  onTrackClick,
}: TrackListProps) => {
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        delay: 250,
        tolerance: 5,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    
    if (!over || !onTrackReorder) return;
    
    const oldIndex = tracks.findIndex((track) => track.id === active.id);
    const newIndex = tracks.findIndex((track) => track.id === over.id);
    
    if (oldIndex !== newIndex) {
      onTrackReorder(oldIndex, newIndex);
    }
  };

  if (!isDraggable) {
    return (
      <List>
        {tracks.map((track, index) => (
          <div key={track.id}>
            <TrackItem
              track={track}
              tracks={tracks}
              index={index}
              onAddToPlaylist={onAddToPlaylist}
              showTrackNumber={showTrackNumber}
              showAlbumLink={showAlbumLink}
              onTrackClick={onTrackClick}
            />
            {index < tracks.length - 1 && <Divider />}
          </div>
        ))}
      </List>
    );
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext
        items={tracks.map(track => track.id)}
        strategy={verticalListSortingStrategy}
      >
        <List>
          {tracks.map((track, index) => (
            <div key={track.id}>
              <SortableTrackItem
                track={track}
                tracks={tracks}
                index={index}
                onAddToPlaylist={onAddToPlaylist}
                showTrackNumber={showTrackNumber}
                showAlbumLink={showAlbumLink}
                onTrackClick={onTrackClick}
              />
              {index < tracks.length - 1 && <Divider />}
            </div>
          ))}
        </List>
      </SortableContext>
    </DndContext>
  );
};

export default TrackList; 