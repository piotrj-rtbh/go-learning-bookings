function roomAvailability(room_id, csrf_token) {
  document.getElementById('check-availability-button').addEventListener("click", function () {
    let html = /*html*/`
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="form-row">
              <div class="col">
                <div class="form-row" id="reservation-dates-modal">
                  <div class="col">
                    <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                  </div>
                  <div class="col">
                    <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                  </div>
                </div>
              </div>
            </div>
        </form>
        `;
    attention.custom({
      msg: html,
      title: "Choose your dates",
      willOpen: () => {
        // we want to display pop-up calendar
        const elem = document.getElementById('reservation-dates-modal');
        const rp = new DateRangePicker(elem, {
          format: 'yyyy-mm-dd',
          showOnFocus: true,
          minDate: new Date(),
        })
      },
      didOpen: () => {
        // we have to remove disabled properties
        // disabled were introduced so that the inputs don't have focus
        document.getElementById("start").removeAttribute("disabled");
        document.getElementById("end").removeAttribute("disabled");
      },
      callback: function (result) {
        console.log("called")

        const form = document.getElementById("check-availability-form");
        const formData = new FormData(form); // create form data of this form
        // we have to add the csrf_token here because it's not in the above form (we could add it there but there's another way)
        formData.append("csrf_token", csrf_token);
        formData.append("room_id", room_id);

        // perform JSON request with POST (we change the method in routes as well)
        fetch('/search-availability-json', {
          method: "post",
          body: formData,
        })
          .then(response => response.json())
          .then(data => {
            if (data.ok) {
              const href = `/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}`
              attention.custom({
                icon: 'success',
                showConfirmButton: false,
                msg: /*html*/`<p>Room is available</p>
                    <p><a href="${href}" class="btn btn-primary">Book now!</a></p>`
              })
            } else {
              attention.error({
                msg: "No availability",
              });
            }
          });
      }
    });
  });
}