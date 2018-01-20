"""Remove uniqueness in Repo

Revision ID: 51d493c4d3e1
Revises: 5ac5404bfcd9
Create Date: 2015-05-11 18:55:46.065354

"""

# revision identifiers, used by Alembic.
revision = '51d493c4d3e1'
down_revision = '5ac5404bfcd9'

from alembic import op
import sqlalchemy as sa


def upgrade():
    ### commands auto generated by Alembic - please adjust! ###
    op.drop_index(u'ix_RepositoryApps_url', table_name='RepositoryApps')
    op.create_index(u'ix_RepositoryApps_url', 'RepositoryApps', ['url'], unique=False)
    ### end Alembic commands ###


def downgrade():
    ### commands auto generated by Alembic - please adjust! ###
    op.drop_index(u'ix_RepositoryApps_url', table_name='RepositoryApps')
    op.create_index(u'ix_RepositoryApps_url', 'RepositoryApps', [u'url'], unique=True)
    ### end Alembic commands ###