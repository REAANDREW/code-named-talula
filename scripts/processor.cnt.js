$recvSync( function ( msg ) {
    try{
        var obj = JSON.parse( msg );

        if(obj.command == 'transform') {
            return transform(obj);
        }
    }catch(err){
        $print(err);
        return "{'error':'transformation error occurred'}"
    }
} );

function transform(obj){
  var output = this['transform_'+obj.id](obj.data)
  return JSON.stringify( output );
}
