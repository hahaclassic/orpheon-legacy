package track_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type TrackSegmentController struct {
	segmentService track.TrackSegmentService
}

func NewTrackSegmentController(segmentService track.TrackSegmentService) *TrackSegmentController {
	return &TrackSegmentController{
		segmentService: segmentService,
	}
}

func (c *TrackSegmentController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Show Segment Stats",
			Run:  c.ShowSegmentStats,
		},
	}
}

func (c *TrackSegmentController) ShowSegmentStats(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter track ID: ")
	scanner.Scan()
	trackID := scanner.Text()

	id, err := uuid.Parse(trackID)
	if err != nil {
		return err
	}

	segments, err := c.segmentService.GetSegments(ctx, id)
	if err != nil {
		return err
	}

	// testTrackID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	// segments := []*entity.Segment{
	// 	{TrackID: testTrackID, Idx: 0, StreamCount: 847, Range: &entity.Range{Start: 0, End: 3}},
	// 	{TrackID: testTrackID, Idx: 1, StreamCount: 923, Range: &entity.Range{Start: 3, End: 6}},
	// 	{TrackID: testTrackID, Idx: 2, StreamCount: 756, Range: &entity.Range{Start: 6, End: 9}},
	// 	{TrackID: testTrackID, Idx: 3, StreamCount: 1024, Range: &entity.Range{Start: 9, End: 12}},
	// 	{TrackID: testTrackID, Idx: 4, StreamCount: 891, Range: &entity.Range{Start: 12, End: 15}},
	// 	{TrackID: testTrackID, Idx: 5, StreamCount: 678, Range: &entity.Range{Start: 15, End: 18}},
	// 	{TrackID: testTrackID, Idx: 6, StreamCount: 945, Range: &entity.Range{Start: 18, End: 21}},
	// 	{TrackID: testTrackID, Idx: 7, StreamCount: 812, Range: &entity.Range{Start: 21, End: 24}},
	// 	{TrackID: testTrackID, Idx: 8, StreamCount: 789, Range: &entity.Range{Start: 24, End: 27}},
	// 	{TrackID: testTrackID, Idx: 9, StreamCount: 876, Range: &entity.Range{Start: 27, End: 30}},
	// 	{TrackID: testTrackID, Idx: 10, StreamCount: 934, Range: &entity.Range{Start: 30, End: 33}},
	// 	{TrackID: testTrackID, Idx: 11, StreamCount: 765, Range: &entity.Range{Start: 33, End: 36}},
	// 	{TrackID: testTrackID, Idx: 12, StreamCount: 823, Range: &entity.Range{Start: 36, End: 39}},
	// 	{TrackID: testTrackID, Idx: 13, StreamCount: 901, Range: &entity.Range{Start: 39, End: 42}},
	// 	{TrackID: testTrackID, Idx: 14, StreamCount: 987, Range: &entity.Range{Start: 42, End: 45}},
	// 	{TrackID: testTrackID, Idx: 15, StreamCount: 654, Range: &entity.Range{Start: 45, End: 48}},
	// 	{TrackID: testTrackID, Idx: 16, StreamCount: 789, Range: &entity.Range{Start: 48, End: 51}},
	// 	{TrackID: testTrackID, Idx: 17, StreamCount: 876, Range: &entity.Range{Start: 51, End: 54}},
	// 	{TrackID: testTrackID, Idx: 18, StreamCount: 945, Range: &entity.Range{Start: 54, End: 57}},
	// 	{TrackID: testTrackID, Idx: 19, StreamCount: 812, Range: &entity.Range{Start: 57, End: 60}},
	// 	{TrackID: testTrackID, Idx: 20, StreamCount: 678, Range: &entity.Range{Start: 60, End: 63}},
	// 	{TrackID: testTrackID, Idx: 21, StreamCount: 923, Range: &entity.Range{Start: 63, End: 66}},
	// 	{TrackID: testTrackID, Idx: 22, StreamCount: 847, Range: &entity.Range{Start: 66, End: 69}},
	// 	{TrackID: testTrackID, Idx: 23, StreamCount: 901, Range: &entity.Range{Start: 69, End: 72}},
	// 	{TrackID: testTrackID, Idx: 24, StreamCount: 765, Range: &entity.Range{Start: 72, End: 75}},
	// 	{TrackID: testTrackID, Idx: 25, StreamCount: 834, Range: &entity.Range{Start: 75, End: 78}},
	// 	{TrackID: testTrackID, Idx: 26, StreamCount: 912, Range: &entity.Range{Start: 78, End: 81}},
	// 	{TrackID: testTrackID, Idx: 27, StreamCount: 789, Range: &entity.Range{Start: 81, End: 84}},
	// 	{TrackID: testTrackID, Idx: 28, StreamCount: 856, Range: &entity.Range{Start: 84, End: 87}},
	// 	{TrackID: testTrackID, Idx: 29, StreamCount: 923, Range: &entity.Range{Start: 87, End: 90}},
	// 	{TrackID: testTrackID, Idx: 30, StreamCount: 847, Range: &entity.Range{Start: 90, End: 93}},
	// 	{TrackID: testTrackID, Idx: 31, StreamCount: 923, Range: &entity.Range{Start: 93, End: 96}},
	// 	{TrackID: testTrackID, Idx: 32, StreamCount: 756, Range: &entity.Range{Start: 96, End: 99}},
	// 	{TrackID: testTrackID, Idx: 33, StreamCount: 1024, Range: &entity.Range{Start: 99, End: 102}},
	// 	{TrackID: testTrackID, Idx: 34, StreamCount: 891, Range: &entity.Range{Start: 102, End: 105}},
	// 	{TrackID: testTrackID, Idx: 35, StreamCount: 678, Range: &entity.Range{Start: 105, End: 108}},
	// 	{TrackID: testTrackID, Idx: 36, StreamCount: 945, Range: &entity.Range{Start: 108, End: 111}},
	// 	{TrackID: testTrackID, Idx: 37, StreamCount: 812, Range: &entity.Range{Start: 111, End: 114}},
	// 	{TrackID: testTrackID, Idx: 38, StreamCount: 789, Range: &entity.Range{Start: 114, End: 117}},
	// 	{TrackID: testTrackID, Idx: 39, StreamCount: 876, Range: &entity.Range{Start: 117, End: 120}},
	// 	{TrackID: testTrackID, Idx: 40, StreamCount: 934, Range: &entity.Range{Start: 120, End: 123}},
	// 	{TrackID: testTrackID, Idx: 41, StreamCount: 765, Range: &entity.Range{Start: 123, End: 126}},
	// 	{TrackID: testTrackID, Idx: 42, StreamCount: 823, Range: &entity.Range{Start: 126, End: 129}},
	// 	{TrackID: testTrackID, Idx: 43, StreamCount: 901, Range: &entity.Range{Start: 129, End: 132}},
	// 	{TrackID: testTrackID, Idx: 44, StreamCount: 987, Range: &entity.Range{Start: 132, End: 135}},
	// 	{TrackID: testTrackID, Idx: 45, StreamCount: 654, Range: &entity.Range{Start: 135, End: 138}},
	// 	{TrackID: testTrackID, Idx: 46, StreamCount: 789, Range: &entity.Range{Start: 138, End: 141}},
	// 	{TrackID: testTrackID, Idx: 47, StreamCount: 876, Range: &entity.Range{Start: 141, End: 144}},
	// 	{TrackID: testTrackID, Idx: 48, StreamCount: 945, Range: &entity.Range{Start: 144, End: 147}},
	// 	{TrackID: testTrackID, Idx: 49, StreamCount: 812, Range: &entity.Range{Start: 147, End: 150}},
	// 	{TrackID: testTrackID, Idx: 50, StreamCount: 678, Range: &entity.Range{Start: 150, End: 153}},
	// 	{TrackID: testTrackID, Idx: 51, StreamCount: 923, Range: &entity.Range{Start: 153, End: 156}},
	// 	{TrackID: testTrackID, Idx: 52, StreamCount: 847, Range: &entity.Range{Start: 156, End: 159}},
	// 	{TrackID: testTrackID, Idx: 53, StreamCount: 901, Range: &entity.Range{Start: 159, End: 162}},
	// 	{TrackID: testTrackID, Idx: 54, StreamCount: 765, Range: &entity.Range{Start: 162, End: 165}},
	// 	{TrackID: testTrackID, Idx: 55, StreamCount: 834, Range: &entity.Range{Start: 165, End: 168}},
	// 	{TrackID: testTrackID, Idx: 56, StreamCount: 912, Range: &entity.Range{Start: 168, End: 171}},
	// 	{TrackID: testTrackID, Idx: 57, StreamCount: 789, Range: &entity.Range{Start: 171, End: 174}},
	// 	{TrackID: testTrackID, Idx: 58, StreamCount: 856, Range: &entity.Range{Start: 174, End: 177}},
	// 	{TrackID: testTrackID, Idx: 59, StreamCount: 923, Range: &entity.Range{Start: 177, End: 180}},
	// }
	output.PrintStatsGraph(segments)
	return nil
}
