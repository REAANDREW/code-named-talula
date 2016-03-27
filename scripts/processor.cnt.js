$recvSync( function ( msg ) {
  var obj = JSON.parse( msg );
  var output = {
    name: obj.firstname + " " + obj.lastname,
    age: obj.age
  };
  return JSON.stringify( output );
} );
