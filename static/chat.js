$(function() {
	var chatLongPoll = function() {
		console.log('test');
		$.get("chat")
		.done(function(data) {
			$("#chatlines").append('<p>' + data + '</p>');
		})
		.always(function() {
			chatLongPoll();
		});
	};

	var sendChat = function() {
		$.post('chatentry', $('#chatbox').val());
		$('#chatbox').val('');
	}

	$('#chatbox').keypress(function(e) {
		if(e.which == 13) {
			sendChat();
		}
	});

	$('#chatenter').click(function() {
		sendChat();
	});

	chatLongPoll();
});