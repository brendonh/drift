CubeSet = function() {
    this.cubes = [];
};

CubeSet.prototype.addCube = function(x, y, color) {
    this.cubes.append([x, y, {'color': color}]);
}

CubeSet.prototype.toGeometry = function() {
    var geom = new THREE.Geometry();
    geom.materials = [];

    var dim = 200, dim_half = dim / 2;

    for (var x = -1; x < 2; x++) {
        this.buildPlane(geom, x, 0, 0, 'z', 'y', 'x', - 1, - 1, dim,  1 ); // px
	    this.buildPlane(geom, x, 0, 0, 'z', 'y', 'x',   1, - 1, dim, -1 ); // nx
	    this.buildPlane(geom, x, 0, 0, 'x', 'z', 'y',   1,   1, dim,  1 ); // py
	    this.buildPlane(geom, x, 0, 0, 'x', 'z', 'y',   1, - 1, dim, -1 ); // ny
	    this.buildPlane(geom, x, 0, 0, 'x', 'y', 'z',   1, - 1, dim,  1 ); // pz
	    this.buildPlane(geom, x, 0, 0, 'x', 'y', 'z', - 1, - 1, dim, -1 ); // nz
    }

    console.log("Faces:", geom.faces.length);

	//geom.computeCentroids();
	geom.mergeVertices();

    return geom;

}

CubeSet.prototype.buildPlane = function(geom, x, y, z, u, v, w, udir, vdir, dim, depth ) {
	var dim_half = dim / 2;
    var offset = geom.vertices.length;

	for ( iy = 0; iy < 2; iy ++ ) {
		for ( ix = 0; ix < 2; ix ++ ) {
			var vector = new THREE.Vector3(dim * x, dim * y, dim * z);
			vector[ u ] += ( ix * dim - dim_half ) * udir;
			vector[ v ] += ( iy * dim - dim_half ) * vdir;
			vector[ w ] += dim_half * depth;
			geom.vertices.push( vector );
		}
	}

	var normal = new THREE.Vector3();
	normal[ w ] = depth > 0 ? 1 : - 1;

    var face = new THREE.Face4( 0 + offset, 2 + offset, 3 + offset, 1 + offset );

	face.normal.copy( normal );
    face.vertexNormals.push( normal, normal, normal, normal );
    
	geom.faces.push( face );
	geom.faceVertexUvs[ 0 ].push( [
		new THREE.UV( 0, 0 ),
		new THREE.UV( 0, 1 ),
		new THREE.UV( 1, 1 ),
		new THREE.UV( 1, 0 )
	] );

}
