function loadPage(url, pop) {

  $.ajax({url: url, success: function(response){
    $('#content').html(jQuery(response).find('#content').html());

    if (!pop) {
      window.history.pushState({url: url}, "", url);
    }
  }});
}

function initModal() {

  $("#theModal").on("show.bs.modal", function(e) {
    var link = $(e.relatedTarget);

    $(this).find(".modal-title").html("<h4 class='modal-title'>Loading...</h4>");
    $(this).find(".modal-body").html("</br><div class='progress'><div class='progress-bar progress-bar-striped active' role='progressbar' style='width:100%'></div></div>")
    $(this).find(".modal-footer").addClass('hide');

    $.get(link.attr("href"), function( response ) {
      $("#theModal").find(".modal-title").html(jQuery(response).filter('.modal-title').html());
      $("#theModal").find(".modal-body").html(jQuery(response).filter('.modal-body').html());
      $("#theModal").find(".modal-footer").html(jQuery(response).filter('.modal-footer').html()).removeClass('hide');
    }, "html");

  });

  $("#theModal").on("hide.bs.modal", function(e) {
    $("body").removeClass("loading");
  });

}

function initMenu() {
  $(document).on({
      ajaxStart: function() { $("body").addClass("loading");    },
       ajaxStop: function() { $("body").removeClass("loading"); }
  });

  $('#menu ul').hide();
  $('#menu ul').children('.current').parent().show();

  $('#menu li a').click(
    function(e) {

      var $parent = $(this).parent();
      var $child = $(this).next();

      if (!$child.is('ul')) {
        if (!$parent.hasClass('active')) {
            $('#menu li').removeClass('active');
            $parent.addClass('active');

            loadPage($(this).attr('href'), false);
            return false;
        } else {
          return false;
        }
      } else {
        if ((!$child.is(':visible'))) {
          $child.slideDown('normal');
          return false;
        } else if (($child.is(':visible'))) {
          return false;
        }
      }

      $('#menu ul:visible').slideUp('normal');

      e.preventDefault();

      }
    );
}

$(document).ready(function() {
  initMenu();
  initModal();

  window.onpopstate = function(e) {
      loadPage(e.state ? e.state.url : null, true);
  };
});



