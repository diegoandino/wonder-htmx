<script>

const friends = async () => { await fetch("").json() }
const currentUserCtx = async () => { await fetch("").json() }

</script>

<h1>Home</h1>
<div id="me">
    <h2>Current User</h2>
    <img source="" alt="profile_picture"/>
    <h3>Song:</h3>
    <img class="current_album_img" source={currentUserCtx.currentAlbumImg} alt="album_img"/>
    <h3 class="current_song">{currentUserCtx.currentSong}</h3>
    <h3 class="current_album_name">{currentUserCtx.currentAlbumName}</h3>
</div>
<div id="friends">
    <h2>Friends</h2>
    <ul id="friend_list">{friends}</ul>
</div>
