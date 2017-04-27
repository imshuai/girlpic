function get_next() {
        $.getJSON("/review/next",
        function (data, textStatus, jqXHR) {
            if (textStatus == "success") {
                var pic = $('#pic');
                pic.empty();
                var a = $('<a></a>');
                var img = $('<img>');
                img.attr('src',data.url);
                img.addClass('img-responsive');
                a.append(img);
                a.attr('href',data.url);
                a.attr('target','_blank');
                pic.append(a);
                pic.data('id',data.id);
            }
        }
    );
}
function get_next_pics(page) {
    $.getJSON("/page/"+page,
        function (data) {
            $.each(data, function (indexInArray, valueOfElement) { 
                var img = $('<img></img>');
                img.attr('src',valueOfElement.url);
                img.addClass('img-responsive center-block');
                var a = $('<a></a>');
                a.attr('href','/detail/'+valueOfElement.id);
                a.attr('target','_blank');
                a.append(img);
                var like = $('<a>');
                var unlike = $('<a>');
                like.addClass('a-like');
                like.attr('href','#');
                like.text('like ');
                $('<span>').addClass('badge').text(valueOfElement.like).appendTo(like);
                unlike.addClass('a-unlike');
                unlike.attr('href','#');
                unlike.text('unlike ');
                $('<span>').addClass('badge').text(valueOfElement.unlike).appendTo(unlike);
                var reviewbox = $('<div>');
                reviewbox.addClass('reviewbox collapse');
                reviewbox.append(like);
                reviewbox.append(unlike);
                var pic = $('<div></div>');
                pic.addClass('thumbnil'); 
                pic.css('margin-top','5px');
                pic.data('id',valueOfElement.id);
                pic.append(a);
                var pics = $('#pics');
                pics.append(pic);
            });
            $('#pics').ready(function(){
                $('#pics').mpmansory({
                    columnClasses: '',
                    breakpoints:{
                        lg: 3, 
                        md: 3, 
                        sm: 6,
                        xs: 12
                    }
                });
            });
        }
    );
    current+=1;
}