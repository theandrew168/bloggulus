package web

import (
	"io"
	"net/http"
	"strconv"

	"github.com/theandrew168/bloggulus/backend/jobqueue"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleJobsList(queue *jobqueue.Queue) http.Handler {
	tmpl := page.NewJobs()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobs, err := queue.ListDead(r.Context())
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}
		data := page.JobsData{
			BaseData: util.GetTemplateBaseData(r, w),

			Jobs: jobs,
		}
		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleJobDeleteForm(queue *jobqueue.Queue) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID, err := strconv.ParseInt(r.PathValue("jobID"), 10, 64)
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		err = queue.Delete(r.Context(), jobID)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		// Redirect back to the jobs page.
		http.Redirect(w, r, "/jobs", http.StatusSeeOther)
	})
}
