$.ajax(
    {
        url: "/basket",
        success: function(result) {
            $(".basket").fadeIn(
                "slow",
                function () {
                    $(".basket").html(result);
                }
            );
        }
    }
);
