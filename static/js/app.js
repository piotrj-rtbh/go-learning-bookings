function Prompt() {
  let toast = function (c) {
      const {
          msg = "",
          icon = "success",
          position = "top-end"
      } = c;

      const Toast = Swal.mixin({
          toast: true,
          title: msg,
          position: position,
          icon: icon,
          showConfirmButton: false,
          timer: 3000,
          timerProgressBar: true,
          didOpen: (toast) => {
              toast.addEventListener('mouseenter', Swal.stopTimer)
              toast.addEventListener('mouseleave', Swal.resumeTimer)
          }
      });

      Toast.fire({});
  }

  let success = function (c) {

      const {
          msg = "",
          title = "",
          footer = "",
      } = c;

      Swal.fire({
          icon: 'success',
          title: title,
          text: msg,
          footer: footer,
      })
  }
  let error = function (c) {
      const {
          msg = "",
          title = "",
          footer = "",
      } = c;

      Swal.fire({
          icon: 'error',
          title: title,
          text: msg,
          footer: footer,
      })
  }

  return {
      toast: toast,
      success: success,
      error: error,
      custom: custom,
  }
}

async function custom(c) {
  const {
      icon = "",
      msg = "",
      title = "",
      showConfirmButton = true
  } = c;

  const { value: result } = await Swal.fire({
      icon: icon,
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      showConfirmButton: showConfirmButton,
      // changing willOpen to a more general - the body will be sent as param
      willOpen: () => {
          if (c.willOpen !== undefined) {
              c.willOpen();
          }
      },
      // same here
      didOpen: () => {
          if (c.didOpen !== undefined) {
              c.didOpen();
          }
      },
      preConfirm: () => {
          return [
              document.getElementById('start').value,
              document.getElementById('end').value
          ]
      }
  })

  // if we have the result
  if (result) {
      // and the result is not because they've clicked the cancel button on the window 
      if (result.dismiss !== Swal.DismissReason.cancel) {
          // and the result is not empty
          if (result.value !== "") {
              // call the callback
              if (c.callback !== undefined) {
                  c.callback(result);
              }
          } else {
              c.callback(false);
          }
      } else {
          c.callback(false);
      }
  }
}

